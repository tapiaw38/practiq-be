package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListTeachersUsecase interface {
		Execute(context.Context, string) (*ListTeachersOutput, apperrors.ApplicationError)
	}

	listTeachersUsecase struct {
		contextFactory appcontext.Factory
	}

	ListTeachersOutput struct {
		Data []UserData `json:"data"`
	}
)

func NewListTeachersUsecase(contextFactory appcontext.Factory) ListTeachersUsecase {
	return &listTeachersUsecase{contextFactory: contextFactory}
}

func (u *listTeachersUsecase) Execute(ctx context.Context, studentID string) (*ListTeachersOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	users, err := app.Repositories.TeacherStudentAssignment.ListTeachers(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
	}
	data := make([]UserData, 0, len(users))
	for _, user := range users {
		data = append(data, toUserData(user))
	}
	return &ListTeachersOutput{Data: data}, nil
}
