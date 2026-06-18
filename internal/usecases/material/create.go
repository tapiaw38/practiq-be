package material

import (
	"context"

	materialRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/material"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, bool, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		CourseID      string
		TeacherID     string
		Title         string `json:"title"`
		Type          string `json:"type"`
		ExtractedText string `json:"extracted_text"`
	}

	CreateOutput struct {
		Data MaterialData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := requesterCanWriteCourse(ctx, app, requesterID, isAdmin, input.CourseID); appErr != nil {
		return nil, appErr
	}

	id, err := app.Repositories.Material.Create(ctx, domain.Material{
		CourseID:      input.CourseID,
		TeacherID:     input.TeacherID,
		Title:         input.Title,
		Type:          input.Type,
		ExtractedText: input.ExtractedText,
		Status:        "uploaded",
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialCreateError, err)
	}

	materials, err := app.Repositories.Material.List(ctx, materialRepo.ListFilter{CourseID: input.CourseID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialListError, err)
	}
	for _, m := range materials {
		if m.ID == id {
			return &CreateOutput{Data: toMaterialData(m)}, nil
		}
	}

	return nil, apperrors.NewInternalError(nil)
}
