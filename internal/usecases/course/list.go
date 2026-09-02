package course

import (
	"context"

	reposCourse "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, ListInput) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListInput struct {
		TeacherID string
		StudentID string
	}

	ListOutput struct {
		Data []CourseData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, input ListInput) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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

	return &ListOutput{Data: data}, nil
}
