package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []LearningStrategyData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	strategies, err := app.Repositories.LearningStrategy.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyListError, err)
	}

	data := make([]LearningStrategyData, 0, len(strategies))
	for _, s := range strategies {
		data = append(data, toStrategyData(s))
	}

	return &ListOutput{Data: data}, nil
}
