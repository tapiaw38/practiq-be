package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListStudentsUsecase interface {
	Execute(context.Context, string) (*UsersOutput, apperrors.ApplicationError)
}

type listStudentsUsecase struct {
	factory appcontext.Factory
}

func NewListStudentsUsecase(factory appcontext.Factory) ListStudentsUsecase {
	return &listStudentsUsecase{factory: factory}
}

func (u *listStudentsUsecase) Execute(ctx context.Context, teacherID string) (*UsersOutput, apperrors.ApplicationError) {
	app := u.factory()
	users, err := app.Repositories.TeacherStudentAssignment.ListStudents(ctx, teacherID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
	}
	data := make([]UserData, 0, len(users))
	for _, user := range users {
		data = append(data, toUserData(user))
	}
	return &UsersOutput{Data: data}, nil
}
