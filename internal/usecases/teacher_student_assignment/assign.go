package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type AssignUsecase interface {
	Execute(context.Context, string, string) (*MutationOutput, apperrors.ApplicationError)
}

type assignUsecase struct {
	factory appcontext.Factory
}

func NewAssignUsecase(factory appcontext.Factory) AssignUsecase {
	return &assignUsecase{factory: factory}
}

func (u *assignUsecase) Execute(ctx context.Context, teacherID, studentID string) (*MutationOutput, apperrors.ApplicationError) {
	app := u.factory()
	if err := app.Repositories.TeacherStudentAssignment.Assign(ctx, domain.TeacherStudentAssignment{
		TeacherID: teacherID,
		StudentID: studentID,
		Status:    "active",
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentCreateError, err)
	}
	return &MutationOutput{Message: "student assigned successfully"}, nil
}
