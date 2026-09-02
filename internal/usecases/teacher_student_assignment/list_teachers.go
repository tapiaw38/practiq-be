package teacherstudentassignment

import (
	"context"

	tsaRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/teacher_student_assignment"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	ListTeachersUsecase interface {
		Execute(context.Context, ListTeachersInput) (*ListTeachersOutput, apperrors.ApplicationError)
	}

	listTeachersUsecase struct {
		contextFactory appcontext.Factory
	}

	ListTeachersInput struct {
		StudentID   string
		Limit       int
		Offset      int
		BearerToken string
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

	ids := make([]string, 0, len(users))
	for _, user := range users {
		ids = append(ids, user.ID)
	}
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, ids)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	data := make([]UserData, 0, len(users))
	for _, user := range users {
		info := names[user.ID]
		data = append(data, toUserData(user, identity.FullName(info, user.ID), info.Email))
	}
	return &ListTeachersOutput{Data: data}, nil
}
