package topic

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Title       string
		Description string
		OrderIndex  int
	}

	UpdateOutput struct {
		Data TopicData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Verify topic exists and check ownership
	topic, err := app.Repositories.Topic.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.TopicGetError, err)
	}
	if topic == nil {
		return nil, apperrors.NewNotFoundError("topic not found")
	}

	if !isSuperAdmin {
		course, err := app.Repositories.Course.Get(ctx, topic.CourseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		if course == nil {
			return nil, apperrors.NewNotFoundError("course not found")
		}
		if course.TeacherID != requesterID {
			return nil, apperrors.NewForbiddenError()
		}
	}

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

	return &UpdateOutput{Data: toTopicData(*t)}, nil
}
