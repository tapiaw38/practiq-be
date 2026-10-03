package school

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func EnsureStudentCanWork(ctx context.Context, app *appcontext.Context, studentID, courseID string) apperrors.ApplicationError {
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil || course == nil || course.SchoolID == "" {

		if err != nil {
			log.Printf("[school] course lookup failed course_id=%s err=%v", courseID, err)
		}
		return nil
	}

	owner := course.TeacherID
	if !subscription.StudentMayWork(ctx, app, course.SchoolID, owner, studentID) {
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

	return nil
}
