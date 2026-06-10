package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type DeleteUsecase interface {
	Execute(context.Context, string) apperrors.ApplicationError
}

type deleteUsecase struct {
	factory appcontext.Factory
}

func NewDeleteUsecase(factory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{factory: factory}
}

func (u *deleteUsecase) Execute(ctx context.Context, id string) apperrors.ApplicationError {
	app := u.factory()

	existing, err := app.Repositories.LearningStrategy.Get(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}
	if existing == nil {
		return apperrors.NewNotFoundError("learning strategy not found")
	}

	if err := app.Repositories.LearningStrategy.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.LearningStrategyDeleteError, err)
	}

	return nil
}
