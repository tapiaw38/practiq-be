package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	DeleteUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, id string) apperrors.ApplicationError
	}

	deleteUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, id string) apperrors.ApplicationError {
	app := u.contextFactory()

	// Verify practice sheet exists and check ownership
	ps, err := app.Repositories.PracticeSheet.Get(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	// Check course ownership for authorization
	if !isAdmin {
		course, err := app.Repositories.Course.Get(ctx, ps.CourseID)
		if err != nil || course == nil {
			return apperrors.NewForbiddenError()
		}
		if course.TeacherID != requesterID {
			return apperrors.NewForbiddenError()
		}
	}

	if err := app.Repositories.PracticeSheet.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.PracticeSheetDeleteError, err)
	}

	return nil
}
