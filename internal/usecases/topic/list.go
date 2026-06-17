package topic

import (
	"context"

	topicRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/topic"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, ListInput) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListInput struct {
		CourseID string
		Limit    int
		Offset   int
	}

	ListOutput struct {
		Data []TopicData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, input ListInput) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	filter := topicRepo.ListFilter{
		CourseID: input.CourseID,
		Limit:    input.Limit,
		Offset:   input.Offset,
	}

	topics, err := app.Repositories.Topic.List(ctx, filter)
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

	return &ListOutput{Data: data}, nil
}
