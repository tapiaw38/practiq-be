package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ReactivateStudentUsecase interface {
		Execute(ctx context.Context, teacherID, studentID string) apperrors.ApplicationError
	}

	reactivateStudentUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewReactivateStudentUsecase(contextFactory appcontext.Factory) ReactivateStudentUsecase {
	return &reactivateStudentUsecase{contextFactory: contextFactory}
}

func (u *reactivateStudentUsecase) Execute(ctx context.Context, teacherID, studentID string) apperrors.ApplicationError {
	app := u.contextFactory()
	scope, appErr := scopeFor(ctx, app, "", teacherID)
	if appErr != nil {
		return appErr
	}
	if scope.SchoolID == "" {
		return apperrors.NewNotFoundError("no school to reactivate in")
	}
	if scope.Enforced() {
		used, appErr := studentsUsed(ctx, app, scope.SchoolID)
		if appErr != nil {
			return appErr
		}
		if !(domain.TeacherSubscription{Plan: scope.Plan, StudentsUsed: used}).CanAddStudent() {
			return apperrors.NewApplicationError(mappings.StudentLimitReachedError, nil)
		}
	}
	if err := app.Repositories.School.SetMemberActive(ctx, scope.SchoolID, studentID, true); err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return nil
}
