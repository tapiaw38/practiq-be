package exercise

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, UpdateInput) (*ExerciseOutput, apperrors.ApplicationError)
}

type updateUsecase struct {
	factory appcontext.Factory
}

type UpdateInput struct {
	Type          string `json:"type"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Explanation   string `json:"explanation"`
	Difficulty    int    `json:"difficulty"`
}

func NewUpdateUsecase(factory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{factory: factory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*ExerciseOutput, apperrors.ApplicationError) {
	app := u.factory()

	if err := app.Repositories.Exercise.Update(ctx, id, domain.Exercise{
		Type:          input.Type,
		Question:      input.Question,
		CorrectAnswer: input.CorrectAnswer,
		Explanation:   input.Explanation,
		Difficulty:    input.Difficulty,
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

	return &ExerciseOutput{Data: toExerciseData(*e)}, nil
}
