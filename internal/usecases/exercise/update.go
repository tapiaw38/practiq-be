package exercise

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Type          string `json:"type"`
		Question      string `json:"question"`
		CorrectAnswer string `json:"correct_answer"`
		Explanation   string `json:"explanation"`
		Difficulty    int    `json:"difficulty"`
		Metadata      string `json:"metadata"`
	}

	UpdateOutput struct {
		Data ExerciseData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Verify exercise exists and check ownership
	exercise, err := app.Repositories.Exercise.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseListError, err)
	}
	if exercise == nil {
		return nil, apperrors.NewNotFoundError("exercise not found")
	}

	if !isAdmin {
		topic, err := app.Repositories.Topic.Get(ctx, exercise.TopicID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.TopicGetError, err)
		}
		if topic == nil {
			return nil, apperrors.NewNotFoundError("topic not found")
		}

		course, err := app.Repositories.Course.Get(ctx, topic.CourseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		if course == nil {
			return nil, apperrors.NewNotFoundError("course not found")
		}
		if course.TeacherID != requesterID {
			return nil, apperrors.NewForbiddenError()
		}
	}

	if appErr := validateFillBlanks(input.Type, input.Question, input.Metadata); appErr != nil {
		return nil, appErr
	}
	if appErr := validateExerciseMediaURL(app, requesterID, input.Metadata, exercise.Metadata); appErr != nil {
		return nil, appErr
	}

	if err := app.Repositories.Exercise.Update(ctx, id, domain.Exercise{
		Type:          input.Type,
		Question:      input.Question,
		CorrectAnswer: input.CorrectAnswer,
		Explanation:   input.Explanation,
		Difficulty:    input.Difficulty,
		Metadata:      input.Metadata,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseUpdateError, err)
	}

	e, err := app.Repositories.Exercise.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseListError, err)
	}
	if e == nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseNotFoundError, nil)
	}

	return &UpdateOutput{Data: toExerciseData(app, *e)}, nil
}
