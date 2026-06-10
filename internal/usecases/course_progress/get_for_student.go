package courseprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GetForStudentUsecase interface {
	Execute(ctx context.Context, studentID, courseID string) (*CourseProgressOutput, apperrors.ApplicationError)
}

type getForStudentUsecase struct {
	factory appcontext.Factory
}

func NewGetForStudentUsecase(factory appcontext.Factory) GetForStudentUsecase {
	return &getForStudentUsecase{factory: factory}
}

func (u *getForStudentUsecase) Execute(ctx context.Context, studentID, courseID string) (*CourseProgressOutput, apperrors.ApplicationError) {
	app := u.factory()

	progress, err := app.Repositories.CourseProgress.Get(ctx, studentID, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseProgressGetError, err)
	}
	if progress == nil {
		return nil, apperrors.NewNotFoundError("course progress not found")
	}

	return &CourseProgressOutput{Data: toProgressData(*progress)}, nil
}
