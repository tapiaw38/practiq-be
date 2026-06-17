package exercise

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	DeleteUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, id string) apperrors.ApplicationError
	}

	deleteUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, id string) apperrors.ApplicationError {
	app := u.contextFactory()

	// Verify exercise exists and check ownership
	exercise, err := app.Repositories.Exercise.Get(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.ExerciseListError, err)
	}
	if exercise == nil {
		return apperrors.NewNotFoundError("exercise not found")
	}

	if !isAdmin {
		topic, err := app.Repositories.Topic.Get(ctx, exercise.TopicID)
		if err != nil {
			return apperrors.NewApplicationError(mappings.TopicGetError, err)
		}
		if topic == nil {
			return apperrors.NewNotFoundError("topic not found")
		}

		course, err := app.Repositories.Course.Get(ctx, topic.CourseID)
		if err != nil {
			return apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		if course == nil {
			return apperrors.NewNotFoundError("course not found")
		}
		if course.TeacherID != requesterID {
			return apperrors.NewForbiddenError()
		}
	}

	if err := app.Repositories.Exercise.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.ExerciseDeleteError, err)
	}

	return nil
}
