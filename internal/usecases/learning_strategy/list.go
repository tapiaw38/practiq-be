package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListUsecase interface {
	Execute(context.Context) (*LearningStrategyListOutput, apperrors.ApplicationError)
}

type listUsecase struct {
	factory appcontext.Factory
}

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context) (*LearningStrategyListOutput, apperrors.ApplicationError) {
	app := u.factory()

	strategies, err := app.Repositories.LearningStrategy.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyListError, err)
	}

	data := make([]LearningStrategyData, 0, len(strategies))
	for _, s := range strategies {
		data = append(data, toStrategyData(s))
	}

	return &LearningStrategyListOutput{Data: data}, nil
}
