package grade

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	schoolUC "github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
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

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// The route lets a teacher ask; this decides which rows they may touch.
	// Without it, opening these routes beyond the platform superadmin would let
	// any teacher edit another school's grades.
	current, err := app.Repositories.Grade.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if current == nil {
		return nil, apperrors.NewNotFoundError("grade not found")
	}
	if appErr := schoolUC.EnsureAdministers(ctx, app, requesterID, isSuperAdmin, current.SchoolID); appErr != nil {
		return nil, appErr
	}

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
