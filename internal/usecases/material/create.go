package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type CreateUsecase interface {
	Execute(context.Context, CreateInput) (*MaterialOutput, apperrors.ApplicationError)
}

type createUsecase struct {
	factory appcontext.Factory
}

type CreateInput struct {
	CourseID      string
	TeacherID     string
	Title         string `json:"title"`
	Type          string `json:"type"`
	ExtractedText string `json:"extracted_text"`
}

func NewCreateUsecase(factory appcontext.Factory) CreateUsecase {
	return &createUsecase{factory: factory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*MaterialOutput, apperrors.ApplicationError) {
	app := u.factory()

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

	materials, err := app.Repositories.Material.List(ctx, input.CourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialListError, err)
	}
	for _, m := range materials {
		if m.ID == id {
			return &MaterialOutput{Data: toMaterialData(m)}, nil
		}
	}

	return nil, apperrors.NewInternalError(nil)
}
