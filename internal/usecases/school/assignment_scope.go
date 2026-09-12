package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// EnsureCanLinkTeacherStudent refuses a requester who does not administer a
// school both the teacher and the student already belong to.
//
// Teacher-student assignment was superadmin-only, and the usecase itself had
// no school check — the route gate was the only thing standing between a
// school admin and linking people outside their school. Opening the route
// without this would let an admin of school A link a teacher from school B
// to a student from school C.
func EnsureCanLinkTeacherStudent(
	ctx context.Context,
	app *appcontext.Context,
	requesterID string,
	isSuperAdmin bool,
	teacherID, studentID string,
) apperrors.ApplicationError {
	if isSuperAdmin {
		return nil
	}

	adminSchools, err := activeSchoolIDs(ctx, app, requesterID, isSchoolAdmin)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	teacherSchools, err := activeSchoolIDs(ctx, app, teacherID, anyRole)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	studentSchools, err := activeSchoolIDs(ctx, app, studentID, anyRole)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	for schoolID := range adminSchools {
		if teacherSchools[schoolID] && studentSchools[schoolID] {
			return nil
		}
	}
	return apperrors.NewForbiddenError()
}

// EnsureCanViewAssignmentsFor refuses a requester who is neither the target
// user nor administers a school the target belongs to. It reads one side of
// a teacher-student link (a teacher's students, or a student's teachers)
// without requiring the caller to know the other side.
func EnsureCanViewAssignmentsFor(
	ctx context.Context,
	app *appcontext.Context,
	requesterID string,
	isSuperAdmin bool,
	targetUserID string,
) apperrors.ApplicationError {
	if isSuperAdmin || requesterID == targetUserID {
		return nil
	}

	adminSchools, err := activeSchoolIDs(ctx, app, requesterID, isSchoolAdmin)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	targetSchools, err := activeSchoolIDs(ctx, app, targetUserID, anyRole)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	for schoolID := range adminSchools {
		if targetSchools[schoolID] {
			return nil
		}
	}
	return apperrors.NewForbiddenError()
}

func isSchoolAdmin(m domain.SchoolMember) bool { return m.Role == domain.SchoolRoleAdmin }
func anyRole(domain.SchoolMember) bool         { return true }

func activeSchoolIDs(
	ctx context.Context,
	app *appcontext.Context,
	userID string,
	include func(domain.SchoolMember) bool,
) (map[string]bool, error) {
	members, err := app.Repositories.School.ListForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	ids := make(map[string]bool, len(members))
	for _, member := range members {
		if member.Active && include(member) {
			ids[member.SchoolID] = true
		}
	}
	return ids, nil
}
