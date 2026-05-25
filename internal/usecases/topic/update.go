package topic

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, UpdateInput) (*TopicOutput, apperrors.ApplicationError)
}

type updateUsecase struct {
	factory appcontext.Factory
}

type UpdateInput struct {
	Title       string
	Description string
	OrderIndex  int
}

func NewUpdateUsecase(factory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{factory: factory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*TopicOutput, apperrors.ApplicationError) {
	app := u.factory()

	if err := app.Repositories.Topic.Update(ctx, id, domain.Topic{
		Title:       input.Title,
		Description: input.Description,
		OrderIndex:  input.OrderIndex,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.TopicUpdateError, err)
	}

	t, err := app.Repositories.Topic.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TopicGetError, err)
	}
	if t == nil {
		return nil, apperrors.NewApplicationError(mappings.TopicNotFoundError, nil)
	}

	return &TopicOutput{Data: toTopicData(*t)}, nil
}
