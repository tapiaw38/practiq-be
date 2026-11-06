package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

func EnsureCanAddStudent(
	ctx context.Context,
	app *appcontext.Context,
	schoolID, teacherID, studentID string,
) apperrors.ApplicationError {

	scope, appErr := scopeFor(ctx, app, schoolID, teacherID)
	if appErr != nil {
		return appErr
	}
	if !scope.Enforced() {
		return nil
	}

	already, appErr := isActiveMember(ctx, app, scope.SchoolID, studentID)
	if appErr != nil {
		return appErr
	}
	if already {
		return nil
	}

	used, appErr := studentsUsed(ctx, app, scope.SchoolID)
	if appErr != nil {
		return appErr
	}

	if !(domain.TeacherSubscription{Plan: scope.Plan, StudentsUsed: used}).CanAddStudent() {
		return apperrors.NewApplicationError(mappings.StudentLimitReachedError, nil)
	}
	return nil
}
