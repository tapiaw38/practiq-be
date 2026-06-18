package exercise

import (
	"context"

	exerciseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/exercise"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, string, bool, ListInput) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListInput struct {
		TopicID string
		Limit   int
		Offset  int
	}

	ListOutput struct {
		Data []ExerciseData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, input ListInput) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := requesterCanWriteTopic(ctx, app, requesterID, isAdmin, input.TopicID); appErr != nil {
		return nil, appErr
	}

	filter := exerciseRepo.ListFilter{
		TopicID: input.TopicID,
		Limit:   input.Limit,
		Offset:  input.Offset,
	}

	exercises, err := app.Repositories.Exercise.List(ctx, filter)
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
