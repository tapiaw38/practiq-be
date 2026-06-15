package course

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
		TeacherID   string
		GradeID     string `json:"grade_id"`
		SubjectID   string `json:"subject_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Level       string `json:"level"`
		Subject     string `json:"subject"`
	}

	CreateOutput struct {
		Data CourseData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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

	defaultStrategy, err := app.Repositories.LearningStrategy.GetByCode(ctx, "kumon")
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}
	if defaultStrategy != nil {
		if _, err := app.Repositories.LearningStrategy.AssignToCourse(ctx, domain.CourseLearningStrategy{
			CourseID:   id,
			StrategyID: defaultStrategy.ID,
			IsDefault:  true,
			Config:     "{}",
		}); err != nil {
			return nil, apperrors.NewApplicationError(mappings.LearningStrategyAssignError, err)
		}
	}

	c, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}

	return &CreateOutput{Data: toCourseData(*c)}, nil
}
