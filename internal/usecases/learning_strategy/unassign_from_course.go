package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	UnassignFromCourseUsecase interface {
		Execute(ctx context.Context, requesterID, id string, isSuperAdmin bool) apperrors.ApplicationError
	}

	unassignFromCourseUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewUnassignFromCourseUsecase(contextFactory appcontext.Factory) UnassignFromCourseUsecase {
	return &unassignFromCourseUsecase{contextFactory: contextFactory}
}

func (u *unassignFromCourseUsecase) Execute(ctx context.Context, requesterID, id string, isSuperAdmin bool) apperrors.ApplicationError {
	app := u.contextFactory()

	existing, err := app.Repositories.LearningStrategy.GetCourseStrategy(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}
	if existing == nil {
		return apperrors.NewNotFoundError("course learning strategy assignment not found")
	}

	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, existing.CourseID); appErr != nil {
		return appErr
	}

	if err := app.Repositories.LearningStrategy.UnassignFromCourse(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.LearningStrategyUnassignError, err)
	}

	return nil
}
