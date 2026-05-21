package exercise

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type CreateUsecase interface {
	Execute(context.Context, CreateInput) (*ExerciseOutput, apperrors.ApplicationError)
}

type createUsecase struct {
	factory appcontext.Factory
}

type CreateInput struct {
	TopicID       string
	Type          string `json:"type"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Explanation   string `json:"explanation"`
	Difficulty    int    `json:"difficulty"`
	Metadata      string `json:"metadata"`
}

func NewCreateUsecase(factory appcontext.Factory) CreateUsecase {
	return &createUsecase{factory: factory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*ExerciseOutput, apperrors.ApplicationError) {
	app := u.factory()

	difficulty := input.Difficulty
	if difficulty < 1 {
		difficulty = 1
	}
	if difficulty > 10 {
		difficulty = 10
	}

	id, err := app.Repositories.Exercise.Create(ctx, domain.Exercise{
		TopicID:       input.TopicID,
		Type:          input.Type,
		Question:      input.Question,
		CorrectAnswer: input.CorrectAnswer,
		Explanation:   input.Explanation,
		Difficulty:    difficulty,
		Metadata:      input.Metadata,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseCreateError, err)
	}

	e, err := app.Repositories.Exercise.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseListError, err)
	}

	return &ExerciseOutput{Data: toExerciseData(*e)}, nil
}
