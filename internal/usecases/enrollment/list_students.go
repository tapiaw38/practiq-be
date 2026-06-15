package enrollment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListStudentsUsecase interface {
		Execute(context.Context, string) (*ListStudentsOutput, apperrors.ApplicationError)
	}

	listStudentsUsecase struct {
		contextFactory appcontext.Factory
	}

	ListStudentsOutput struct {
		Data []StudentData `json:"data"`
	}
)

func NewListStudentsUsecase(contextFactory appcontext.Factory) ListStudentsUsecase {
	return &listStudentsUsecase{contextFactory: contextFactory}
}

func (u *listStudentsUsecase) Execute(ctx context.Context, courseID string) (*ListStudentsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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

	return &ListStudentsOutput{Data: data}, nil
}
