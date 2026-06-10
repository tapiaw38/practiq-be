package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GetUsecase interface {
	Execute(context.Context, string) (*LearningStrategyOutput, apperrors.ApplicationError)
}

type getUsecase struct {
	factory appcontext.Factory
}

func NewGetUsecase(factory appcontext.Factory) GetUsecase {
	return &getUsecase{factory: factory}
}

func (u *getUsecase) Execute(ctx context.Context, id string) (*LearningStrategyOutput, apperrors.ApplicationError) {
	app := u.factory()

	strategy, err := app.Repositories.LearningStrategy.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}
	if strategy == nil {
		return nil, apperrors.NewNotFoundError("learning strategy not found")
	}

	return &LearningStrategyOutput{Data: toStrategyData(*strategy)}, nil
}
