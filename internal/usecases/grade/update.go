package grade

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
		Name        string
		Description string
		VisualTheme string
	}

	UpdateOutput struct {
		Data GradeData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	visualTheme := strings.TrimSpace(input.VisualTheme)
	if visualTheme != "primary" && visualTheme != "secondary" {
		return nil, apperrors.NewBadRequestError("visual_theme must be primary or secondary")
	}

	if err := app.Repositories.Grade.Update(ctx, id, domain.Grade{
		Name:        input.Name,
		Description: input.Description,
		VisualTheme: visualTheme,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeUpdateError, err)
	}

	grade, err := app.Repositories.Grade.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if grade == nil {
		return nil, apperrors.NewApplicationError(mappings.GradeNotFoundError, nil)
	}

	return &UpdateOutput{Data: toGradeData(*grade)}, nil
}
