package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type UnassignUsecase interface {
	Execute(context.Context, string, string) (*MutationOutput, apperrors.ApplicationError)
}

type unassignUsecase struct {
	factory appcontext.Factory
}

func NewUnassignUsecase(factory appcontext.Factory) UnassignUsecase {
	return &unassignUsecase{factory: factory}
}

func (u *unassignUsecase) Execute(ctx context.Context, teacherID, studentID string) (*MutationOutput, apperrors.ApplicationError) {
	app := u.factory()
	if err := app.Repositories.TeacherStudentAssignment.Unassign(ctx, teacherID, studentID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentDeleteError, err)
	}
	return &MutationOutput{Message: "student unassigned successfully"}, nil
}
