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
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, topicID string, in CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Type          string `json:"type" binding:"required"`
		Question      string `json:"question" binding:"required"`
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

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, topicID string, in CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := requesterCanWriteTopic(ctx, app, requesterID, isSuperAdmin, topicID); appErr != nil {
		return nil, appErr
	}

	if appErr := validateFillBlanks(in.Type, in.Question, in.Metadata, in.CorrectAnswer); appErr != nil {
		return nil, appErr
	}
	if appErr := validateExerciseMediaURL(app, requesterID, in.Metadata, ""); appErr != nil {
		return nil, appErr
	}

	difficulty := in.Difficulty
	if difficulty < 1 {
		difficulty = 1
	}
	if difficulty > 10 {
		difficulty = 10
	}

	metadata := storeTeacherImage(ctx, app, requesterID, in.Metadata, domain.Exercise{})

	id, err := app.Repositories.Exercise.Create(ctx, domain.Exercise{
		TopicID:       topicID,
		Type:          in.Type,
		Question:      in.Question,
		CorrectAnswer: in.CorrectAnswer,
		Explanation:   in.Explanation,
		Difficulty:    difficulty,
		Metadata:      metadata,
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
