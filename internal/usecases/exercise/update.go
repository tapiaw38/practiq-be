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
		Execute(context.Context, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
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

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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

	return &UpdateOutput{Data: toExerciseData(*e)}, nil
}
