package school

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

// JoinTeacherSchool adds a student to the school of the teacher who just took
// them on.
//
// Called from every path where that happens, so the memberships the migration
// backfilled keep being created as the product runs. Skipping one would leave
// students who belong to a teacher but to no school, and they would vanish the
// day anything starts filtering by it.
//
// Today it resolves the teacher's personal school. A teacher at an institution
// inviting a student needs the invitation itself to name the school, which
// comes with institutions.
func JoinTeacherSchool(ctx context.Context, app *appcontext.Context, teacherID, studentID string) {
	school, err := app.Repositories.School.GetPersonal(ctx, teacherID)
	if err != nil {
		log.Printf("[schools] school lookup failed teacher_id=%s err=%v", teacherID, err)
		return
	}
	if school == nil {
		return
	}

	// Swallowed on purpose: the student is already linked to the teacher by the
	// caller, and failing the whole operation over a membership row would undo
	// something that worked. The next sync or redemption adds it.
	if err := app.Repositories.School.AddMember(ctx, domain.SchoolMember{
		SchoolID: school.ID,
		UserID:   studentID,
		Role:     domain.SchoolRoleStudent,
	}); err != nil {
		log.Printf("[schools] membership failed school_id=%s user_id=%s err=%v", school.ID, studentID, err)
	}
}
