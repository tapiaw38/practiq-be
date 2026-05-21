package topic

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type CreateUsecase interface {
	Execute(context.Context, CreateInput) (*TopicOutput, apperrors.ApplicationError)
}

type createUsecase struct {
	factory appcontext.Factory
}

type CreateInput struct {
	CourseID    string
	Title       string `json:"title"`
	Description string `json:"description"`
	OrderIndex  int    `json:"order_index"`
}

func NewCreateUsecase(factory appcontext.Factory) CreateUsecase {
	return &createUsecase{factory: factory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*TopicOutput, apperrors.ApplicationError) {
	app := u.factory()

	id, err := app.Repositories.Topic.Create(ctx, domain.Topic{
		CourseID:    input.CourseID,
		Title:       input.Title,
		Description: input.Description,
		OrderIndex:  input.OrderIndex,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TopicCreateError, err)
	}

	topics, err := app.Repositories.Topic.List(ctx, input.CourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TopicListError, err)
	}

	for _, t := range topics {
		if t.ID == id {
			return &TopicOutput{Data: toTopicData(t)}, nil
		}
	}

	return nil, apperrors.NewInternalError(nil)
}
