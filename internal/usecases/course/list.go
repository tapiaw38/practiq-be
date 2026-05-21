package course

import (
	"context"

	reposCourse "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListUsecase interface {
	Execute(context.Context, ListInput) (*CourseListOutput, apperrors.ApplicationError)
}

type listUsecase struct {
	factory appcontext.Factory
}

type ListInput struct {
	TeacherID string
	StudentID string
}

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context, input ListInput) (*CourseListOutput, apperrors.ApplicationError) {
	app := u.factory()

	courses, err := app.Repositories.Course.List(ctx, reposCourse.ListFilterOptions{
		TeacherID: input.TeacherID,
		StudentID: input.StudentID,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseListError, err)
	}

	var data []CourseData
	for _, c := range courses {
		data = append(data, toCourseData(c))
	}
	if data == nil {
		data = []CourseData{}
	}

	return &CourseListOutput{Data: data}, nil
}
