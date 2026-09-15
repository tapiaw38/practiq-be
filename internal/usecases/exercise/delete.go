package exercise

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

	// Verify exercise exists and check ownership
	exercise, err := app.Repositories.Exercise.Get(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.ExerciseListError, err)
	}
	if exercise == nil {
		return apperrors.NewNotFoundError("exercise not found")
	}

	if appErr := requesterCanWriteTopic(ctx, app, requesterID, isSuperAdmin, exercise.TopicID); appErr != nil {
		return appErr
	}

	if err := app.Repositories.Exercise.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.ExerciseDeleteError, err)
	}

	return nil
}
