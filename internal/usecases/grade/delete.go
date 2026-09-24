package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	schoolUC "github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	DeleteUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string) apperrors.ApplicationError
	}

	deleteUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string) apperrors.ApplicationError {
	app := u.contextFactory()

	// The route lets a teacher ask; this decides which rows they may touch.
	// Without it, opening these routes beyond the platform superadmin would let
	// any teacher edit another school's grades.
	current, err := app.Repositories.Grade.Get(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if current == nil {
		return apperrors.NewNotFoundError("grade not found")
	}
	if appErr := schoolUC.EnsureAdministers(ctx, app, requesterID, isSuperAdmin, current.SchoolID); appErr != nil {
		return appErr
	}

	if err := app.Repositories.Grade.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.GradeDeleteError, err)
	}

	return nil
}
