package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListStudentsUsecase interface {
	Execute(context.Context, string) (*StudentsListOutput, apperrors.ApplicationError)
}

type listStudentsUsecase struct {
	factory appcontext.Factory
}

func NewListStudentsUsecase(factory appcontext.Factory) ListStudentsUsecase {
	return &listStudentsUsecase{factory: factory}
}

func (u *listStudentsUsecase) Execute(ctx context.Context, courseID string) (*StudentsListOutput, apperrors.ApplicationError) {
	app := u.factory()

	students, err := app.Repositories.Enrollment.ListStudents(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentListError, err)
	}

	var data []StudentData
	for _, s := range students {
		data = append(data, toStudentData(s))
	}
	if data == nil {
		data = []StudentData{}
	}

	return &StudentsListOutput{Data: data}, nil
}
