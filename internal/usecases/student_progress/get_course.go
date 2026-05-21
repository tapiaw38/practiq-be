package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GetCourseProgressUsecase interface {
	Execute(context.Context, string, string) (*ProgressListOutput, apperrors.ApplicationError)
}

type getCourseProgressUsecase struct {
	factory appcontext.Factory
}

func NewGetCourseProgressUsecase(factory appcontext.Factory) GetCourseProgressUsecase {
	return &getCourseProgressUsecase{factory: factory}
}

func (u *getCourseProgressUsecase) Execute(ctx context.Context, studentID, courseID string) (*ProgressListOutput, apperrors.ApplicationError) {
	app := u.factory()

	list, err := app.Repositories.StudentProgress.ListByStudentAndCourse(ctx, studentID, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	var data []ProgressData
	for _, p := range list {
		data = append(data, toProgressData(p))
	}
	if data == nil {
		data = []ProgressData{}
	}

	return &ProgressListOutput{Data: data}, nil
}
