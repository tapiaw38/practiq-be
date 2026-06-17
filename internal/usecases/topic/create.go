package topic

import (
	"context"

	topicRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/topic"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(context.Context, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		CourseID    string
		Title       string `json:"title"`
		Description string `json:"description"`
		OrderIndex  int    `json:"order_index"`
	}

	CreateOutput struct {
		Data TopicData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	id, err := app.Repositories.Topic.Create(ctx, domain.Topic{
		CourseID:    input.CourseID,
		Title:       input.Title,
		Description: input.Description,
		OrderIndex:  input.OrderIndex,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TopicCreateError, err)
	}

	topics, err := app.Repositories.Topic.List(ctx, topicRepo.ListFilter{CourseID: input.CourseID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TopicListError, err)
	}

	for _, t := range topics {
		if t.ID == id {
			return &CreateOutput{Data: toTopicData(t)}, nil
		}
	}

	return nil, apperrors.NewInternalError(nil)
}
