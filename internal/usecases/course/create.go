package course

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type CreateUsecase interface {
	Execute(context.Context, CreateInput) (*CourseOutput, apperrors.ApplicationError)
}

type createUsecase struct {
	factory appcontext.Factory
}

type CreateInput struct {
	TeacherID   string
	GradeID     string `json:"grade_id"`
	SubjectID   string `json:"subject_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Subject     string `json:"subject"`
}

func NewCreateUsecase(factory appcontext.Factory) CreateUsecase {
	return &createUsecase{factory: factory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CourseOutput, apperrors.ApplicationError) {
	app := u.factory()

	if strings.TrimSpace(input.GradeID) == "" {
		return nil, apperrors.NewBadRequestError("grade_id is required")
	}
	if strings.TrimSpace(input.SubjectID) == "" {
		return nil, apperrors.NewBadRequestError("subject_id is required")
	}

	id, err := app.Repositories.Course.Create(ctx, domain.Course{
		TeacherID:   input.TeacherID,
		GradeID:     input.GradeID,
		SubjectID:   input.SubjectID,
		Title:       input.Title,
		Description: input.Description,
		Level:       input.Level,
		Subject:     input.Subject,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseCreateError, err)
	}

	c, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}

	return &CourseOutput{Data: toCourseData(*c)}, nil
}
