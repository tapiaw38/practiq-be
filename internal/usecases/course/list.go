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
		Execute(context.Context, string, string, string) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []CourseData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, teacherID, studentID, schoolID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	courses, err := app.Repositories.Course.List(ctx, reposCourse.ListFilterOptions{
		TeacherID: teacherID,
		StudentID: studentID,
		SchoolID:  schoolID,
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
