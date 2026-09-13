package learningstrategy

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(context.Context, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
	}

	CreateOutput struct {
		Data LearningStrategyData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if strings.TrimSpace(input.Name) == "" {
		return nil, apperrors.NewBadRequestError("name is required")
	}
	if strings.TrimSpace(input.Code) == "" {
		return nil, apperrors.NewBadRequestError("code is required")
	}

	id, err := app.Repositories.LearningStrategy.Create(ctx, domain.LearningStrategy{
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
		Status:      "active",
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyCreateError, err)
	}

	strategy, err := app.Repositories.LearningStrategy.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}

	return &CreateOutput{Data: toStrategyData(*strategy)}, nil
}
