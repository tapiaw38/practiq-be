package exercise

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, string) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []ExerciseData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, topicID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	exercises, err := app.Repositories.Exercise.List(ctx, topicID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseListError, err)
	}

	var data []ExerciseData
	for _, e := range exercises {
		data = append(data, toExerciseData(e))
	}
	if data == nil {
		data = []ExerciseData{}
	}

	return &ListOutput{Data: data}, nil
}
