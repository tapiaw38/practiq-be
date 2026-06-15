package courseprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetForStudentUsecase interface {
		Execute(ctx context.Context, studentID, courseID string) (*GetForStudentOutput, apperrors.ApplicationError)
	}

	getForStudentUsecase struct {
		contextFactory appcontext.Factory
	}

	GetForStudentOutput struct {
		Data CourseProgressData `json:"data"`
	}
)

func NewGetForStudentUsecase(contextFactory appcontext.Factory) GetForStudentUsecase {
	return &getForStudentUsecase{contextFactory: contextFactory}
}

func (u *getForStudentUsecase) Execute(ctx context.Context, studentID, courseID string) (*GetForStudentOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	progress, err := app.Repositories.CourseProgress.Get(ctx, studentID, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseProgressGetError, err)
	}
	if progress == nil {
		return nil, apperrors.NewNotFoundError("course progress not found")
	}

	return &GetForStudentOutput{Data: toProgressData(*progress)}, nil
}
