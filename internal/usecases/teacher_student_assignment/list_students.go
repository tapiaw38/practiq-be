package teacherstudentassignment

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
		Data []UserData `json:"data"`
	}
)

func NewListStudentsUsecase(contextFactory appcontext.Factory) ListStudentsUsecase {
	return &listStudentsUsecase{contextFactory: contextFactory}
}

func (u *listStudentsUsecase) Execute(ctx context.Context, teacherID string) (*ListStudentsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	users, err := app.Repositories.TeacherStudentAssignment.ListStudents(ctx, teacherID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
	}
	data := make([]UserData, 0, len(users))
	for _, user := range users {
		data = append(data, toUserData(user))
	}
	return &ListStudentsOutput{Data: data}, nil
}
