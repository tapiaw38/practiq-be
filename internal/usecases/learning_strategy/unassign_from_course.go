package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UnassignFromCourseUsecase interface {
		Execute(context.Context, string) apperrors.ApplicationError
	}

	unassignFromCourseUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewUnassignFromCourseUsecase(contextFactory appcontext.Factory) UnassignFromCourseUsecase {
	return &unassignFromCourseUsecase{contextFactory: contextFactory}
}

func (u *unassignFromCourseUsecase) Execute(ctx context.Context, id string) apperrors.ApplicationError {
	app := u.contextFactory()

	existing, err := app.Repositories.LearningStrategy.GetCourseStrategy(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}
	if existing == nil {
		return apperrors.NewNotFoundError("course learning strategy assignment not found")
	}

	if err := app.Repositories.LearningStrategy.UnassignFromCourse(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.LearningStrategyUnassignError, err)
	}

	return nil
}
