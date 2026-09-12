package teacherstudentassignment

import (
	"context"

	tsaRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/teacher_student_assignment"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	ListStudentsUsecase interface {
		Execute(context.Context, string, bool, ListStudentsInput) (*ListStudentsOutput, apperrors.ApplicationError)
	}

	listStudentsUsecase struct {
		contextFactory appcontext.Factory
	}

	ListStudentsInput struct {
		TeacherID   string
		Limit       int
		Offset      int
		BearerToken string
	}

	ListStudentsOutput struct {
		Data []UserData `json:"data"`
	}
)

func NewListStudentsUsecase(contextFactory appcontext.Factory) ListStudentsUsecase {
	return &listStudentsUsecase{contextFactory: contextFactory}
}

func (u *listStudentsUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, input ListStudentsInput) (*ListStudentsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := school.EnsureCanViewAssignmentsFor(ctx, app, requesterID, isSuperAdmin, input.TeacherID); appErr != nil {
		return nil, appErr
	}

	filter := tsaRepo.ListFilter{
		UserID: input.TeacherID,
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	users, err := app.Repositories.TeacherStudentAssignment.ListStudents(ctx, filter)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
	}

	ids := make([]string, 0, len(users))
	for _, user := range users {
		ids = append(ids, user.ID)
	}
	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, ids)
	if appErr != nil {
		return nil, appErr
	}

	data := make([]UserData, 0, len(users))
	for _, user := range users {
		info := names[user.ID]
		data = append(data, toUserData(user, identity.FullName(info, user.ID), info.Email))
	}
	return &ListStudentsOutput{Data: data}, nil
}
