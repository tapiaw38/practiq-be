package course

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

	course, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return apperrors.NewNotFoundError("course not found")
	}

	if !isAdmin && course.TeacherID != requesterID {
		return apperrors.NewForbiddenError()
	}

	if err := app.Repositories.Course.Delete(ctx, id); err != nil {
		return apperrors.NewApplicationError(mappings.CourseDeleteError, err)
	}

	return nil
}
