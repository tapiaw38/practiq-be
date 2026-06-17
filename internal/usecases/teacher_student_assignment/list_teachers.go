package teacherstudentassignment

import (
	"context"

	tsaRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/teacher_student_assignment"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListTeachersUsecase interface {
		Execute(context.Context, ListTeachersInput) (*ListTeachersOutput, apperrors.ApplicationError)
	}

	listTeachersUsecase struct {
		contextFactory appcontext.Factory
	}

	ListTeachersInput struct {
		StudentID string
		Limit     int
		Offset    int
	}

	ListTeachersOutput struct {
		Data []UserData `json:"data"`
	}
)

func NewListTeachersUsecase(contextFactory appcontext.Factory) ListTeachersUsecase {
	return &listTeachersUsecase{contextFactory: contextFactory}
}

func (u *listTeachersUsecase) Execute(ctx context.Context, input ListTeachersInput) (*ListTeachersOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	filter := tsaRepo.ListFilter{
		UserID: input.StudentID,
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	users, err := app.Repositories.TeacherStudentAssignment.ListTeachers(ctx, filter)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
	}
	data := make([]UserData, 0, len(users))
	for _, user := range users {
		data = append(data, toUserData(user))
	}
	return &ListTeachersOutput{Data: data}, nil
}
