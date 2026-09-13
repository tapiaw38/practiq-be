package grade

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	CreateUsecase interface {
		Execute(context.Context, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Name        string
		Description string
		VisualTheme string
		CreatedBy   string
		SchoolID    string
	}

	CreateOutput struct {
		Data GradeData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Resolved here rather than taken from the request: a caller must not be
	// able to file a grade or subject under a school they do not belong to.
	schoolID, appErr := school.OwnedSchoolIDSelected(ctx, app, input.CreatedBy, input.SchoolID)
	if appErr != nil {
		return nil, appErr
	}

	visualTheme := strings.TrimSpace(input.VisualTheme)
	if visualTheme == "" {
		visualTheme = "primary"
	}
	if visualTheme != "primary" && visualTheme != "secondary" {
		return nil, apperrors.NewBadRequestError("visual_theme must be primary or secondary")
	}

	id, err := app.Repositories.Grade.Create(ctx, domain.Grade{
		SchoolID:    schoolID,
		Name:        input.Name,
		Description: input.Description,
		VisualTheme: visualTheme,
		CreatedBy:   input.CreatedBy,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeCreateError, err)
	}

	grade, err := app.Repositories.Grade.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if grade == nil {
		return nil, apperrors.NewApplicationError(mappings.GradeNotFoundError, nil)
	}

	return &CreateOutput{Data: toGradeData(*grade)}, nil
}
