package grade

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/pgerr"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

const gradeNameConstraint = "idx_grades_school_name"

type (
	CreateUsecase interface {
		Execute(ctx context.Context, createdBy, schoolID string, isSuperAdmin bool, in CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		VisualTheme string `json:"visual_theme"`
	}

	CreateOutput struct {
		Data GradeData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, createdBy, requestedSchoolID string, isSuperAdmin bool, in CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	schoolID, appErr := school.OwnedSchoolIDSelected(ctx, app, createdBy, isSuperAdmin, requestedSchoolID)
	if appErr != nil {
		return nil, appErr
	}

	visualTheme := strings.TrimSpace(in.VisualTheme)
	if visualTheme == "" {
		visualTheme = "primary"
	}
	if visualTheme != "primary" && visualTheme != "secondary" {
		return nil, apperrors.NewBadRequestError("visual_theme must be primary or secondary")
	}

	id, err := app.Repositories.Grade.Create(ctx, domain.Grade{
		SchoolID:    schoolID,
		Name:        in.Name,
		Description: in.Description,
		VisualTheme: visualTheme,
		CreatedBy:   createdBy,
	})
	if err != nil {
		if pgerr.IsUniqueViolation(err, gradeNameConstraint) {
			return nil, apperrors.NewConflictError("ya existe un grado con ese nombre en esta escuela")
		}
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
