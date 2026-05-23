package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListTeachersUsecase interface {
	Execute(context.Context, string) (*UsersOutput, apperrors.ApplicationError)
}

type listTeachersUsecase struct {
	factory appcontext.Factory
}

func NewListTeachersUsecase(factory appcontext.Factory) ListTeachersUsecase {
	return &listTeachersUsecase{factory: factory}
}

func (u *listTeachersUsecase) Execute(ctx context.Context, studentID string) (*UsersOutput, apperrors.ApplicationError) {
	app := u.factory()
	users, err := app.Repositories.TeacherStudentAssignment.ListTeachers(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
	}
	data := make([]UserData, 0, len(users))
	for _, user := range users {
		data = append(data, toUserData(user))
	}
	return &UsersOutput{Data: data}, nil
}
