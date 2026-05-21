package topic

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListUsecase interface {
	Execute(context.Context, string) (*TopicListOutput, apperrors.ApplicationError)
}

type listUsecase struct {
	factory appcontext.Factory
}

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context, courseID string) (*TopicListOutput, apperrors.ApplicationError) {
	app := u.factory()

	topics, err := app.Repositories.Topic.List(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TopicListError, err)
	}

	var data []TopicData
	for _, t := range topics {
		data = append(data, toTopicData(t))
	}
	if data == nil {
		data = []TopicData{}
	}

	return &TopicListOutput{Data: data}, nil
}
