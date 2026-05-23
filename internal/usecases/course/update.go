package course

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, UpdateInput) (*CourseOutput, apperrors.ApplicationError)
}

type updateUsecase struct {
	factory appcontext.Factory
}

type UpdateInput struct {
	GradeID     string `json:"grade_id"`
	SubjectID   string `json:"subject_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Subject     string `json:"subject"`
}

func NewUpdateUsecase(factory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{factory: factory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*CourseOutput, apperrors.ApplicationError) {
	app := u.factory()

	if strings.TrimSpace(input.GradeID) == "" {
		return nil, apperrors.NewBadRequestError("grade_id is required")
	}
	if strings.TrimSpace(input.SubjectID) == "" {
		return nil, apperrors.NewBadRequestError("subject_id is required")
	}

	if err := app.Repositories.Course.Update(ctx, id, domain.Course{
		GradeID:     input.GradeID,
		SubjectID:   input.SubjectID,
		Title:       input.Title,
		Description: input.Description,
		Level:       input.Level,
		Subject:     input.Subject,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseUpdateError, err)
	}

	c, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}

	return &CourseOutput{Data: toCourseData(*c)}, nil
}
