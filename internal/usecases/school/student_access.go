package school

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

// EnsureStudentCanWork refuses a student their school deactivated.
//
// Deactivating only ever stopped them counting against the plan. Their
// enrolments were untouched and nothing on the way in ever looked at the
// membership, so a teacher who dropped from fifteen students to five kept ten
// of them working — while the screen told the teacher those ten had lost
// access. This is what makes that sentence true.
//
// Reading is still allowed, deliberately: a deactivated student keeps their
// notebooks, their marks and their history. What stops is adding to them. The
// student did not make this decision and should not lose what they wrote.
func EnsureStudentCanWork(ctx context.Context, app *appcontext.Context, studentID, courseID string) apperrors.ApplicationError {
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil || course == nil || course.SchoolID == "" {
		// A course we cannot place in a school cannot be judged. Refusing here
		// would lock working students out over a failed lookup.
		if err != nil {
			log.Printf("[school] course lookup failed course_id=%s err=%v", courseID, err)
		}
		return nil
	}

	// A school whose trial ran out and that nobody pays for stops taking work
	// from anybody, not just from students over a cap — there is no cap left.
	owner := course.TeacherID
	if !subscription.StudentsCanWork(ctx, app, course.SchoolID, owner) {
		return apperrors.NewForbiddenError()
	}

	memberships, err := app.Repositories.School.ListForUser(ctx, studentID)
	if err != nil {
		log.Printf("[school] membership lookup failed student_id=%s err=%v", studentID, err)
		return nil
	}
	for _, membership := range memberships {
		if membership.SchoolID != course.SchoolID {
			continue
		}
		if !membership.Active {
			return apperrors.NewForbiddenError()
		}
		return nil
	}
	// No membership row at all: a student who predates schools, or one reached
	// through a grade rather than a school. Not somebody's downgrade.
	return nil
}
