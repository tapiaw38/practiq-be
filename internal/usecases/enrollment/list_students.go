package enrollment

import (
	"context"

	enrollmentRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/enrollment"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListStudentsUsecase interface {
		Execute(context.Context, ListStudentsInput) (*ListStudentsOutput, apperrors.ApplicationError)
	}

	listStudentsUsecase struct {
		contextFactory appcontext.Factory
	}

	ListStudentsInput struct {
		CourseID string
		Limit    int
		Offset   int
	}

	ListStudentsOutput struct {
		Data []StudentData `json:"data"`
	}
)

func NewListStudentsUsecase(contextFactory appcontext.Factory) ListStudentsUsecase {
	return &listStudentsUsecase{contextFactory: contextFactory}
}

func (u *listStudentsUsecase) Execute(ctx context.Context, input ListStudentsInput) (*ListStudentsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	filter := enrollmentRepo.ListFilter{
		CourseID: input.CourseID,
		Limit:    input.Limit,
		Offset:   input.Offset,
	}

	students, err := app.Repositories.Enrollment.ListStudents(ctx, filter)
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

	return &ListStudentsOutput{Data: data}, nil
}
