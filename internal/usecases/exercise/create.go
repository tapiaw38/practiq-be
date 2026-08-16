package exercise

import (
	"context"

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
		TopicID       string
		Type          string `json:"type"`
		Question      string `json:"question"`
		CorrectAnswer string `json:"correct_answer"`
		Explanation   string `json:"explanation"`
		Difficulty    int    `json:"difficulty"`
		Metadata      string `json:"metadata"`
	}

	CreateOutput struct {
		Data ExerciseData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := requesterCanWriteTopic(ctx, app, requesterID, isAdmin, input.TopicID); appErr != nil {
		return nil, appErr
	}

	if appErr := validateFillBlanks(input.Type, input.Question, input.Metadata, input.CorrectAnswer); appErr != nil {
		return nil, appErr
	}
	if appErr := validateExerciseMediaURL(app, requesterID, input.Metadata, ""); appErr != nil {
		return nil, appErr
	}

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

	return &CreateOutput{Data: toExerciseData(app, *e)}, nil
}
