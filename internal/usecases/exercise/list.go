package exercise

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListUsecase interface {
	Execute(context.Context, string) (*ExerciseListOutput, apperrors.ApplicationError)
}

type listUsecase struct {
	factory appcontext.Factory
}

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context, topicID string) (*ExerciseListOutput, apperrors.ApplicationError) {
	app := u.factory()

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

	return &ExerciseListOutput{Data: data}, nil
}
