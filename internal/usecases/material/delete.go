package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
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

	// Verify material exists and check ownership
	material, err := app.Repositories.Material.Get(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.MaterialGetError, err)
	}
	if material == nil {
		return apperrors.NewNotFoundError("material not found")
	}

	if !isSuperAdmin && material.TeacherID != requesterID {
		return apperrors.NewForbiddenError()
	}

	if err := app.Repositories.Material.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.MaterialDeleteError, err)
	}

	return nil
}
