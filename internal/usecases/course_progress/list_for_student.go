package courseprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListForStudentUsecase interface {
	Execute(ctx context.Context, studentID string) (*CourseProgressListOutput, apperrors.ApplicationError)
}

type listForStudentUsecase struct {
	factory appcontext.Factory
}

func NewListForStudentUsecase(factory appcontext.Factory) ListForStudentUsecase {
	return &listForStudentUsecase{factory: factory}
}

func (u *listForStudentUsecase) Execute(ctx context.Context, studentID string) (*CourseProgressListOutput, apperrors.ApplicationError) {
	app := u.factory()

	progressList, err := app.Repositories.CourseProgress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseProgressListError, err)
	}

	data := make([]CourseProgressData, 0, len(progressList))
	for _, p := range progressList {
		data = append(data, toProgressData(p))
	}

	return &CourseProgressListOutput{Data: data}, nil
}
