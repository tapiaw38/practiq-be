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
	UpdateUsecase interface {
		Execute(context.Context, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	UpdateOutput struct {
		Data LearningStrategyData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	existing, err := app.Repositories.LearningStrategy.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}
	if existing == nil {
		return nil, apperrors.NewNotFoundError("learning strategy not found")
	}

	name := input.Name
	if strings.TrimSpace(name) == "" {
		name = existing.Name
	}
	code := input.Code
	if strings.TrimSpace(code) == "" {
		code = existing.Code
	}
	status := input.Status
	if strings.TrimSpace(status) == "" {
		status = existing.Status
	}

	if err := app.Repositories.LearningStrategy.Update(ctx, id, domain.LearningStrategy{
		Name:        name,
		Code:        code,
		Description: input.Description,
		Status:      status,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyUpdateError, err)
	}

	strategy, err := app.Repositories.LearningStrategy.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}

	return &UpdateOutput{Data: toStrategyData(*strategy)}, nil
}
