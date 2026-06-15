package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	AssignUsecase interface {
		Execute(context.Context, string, string) (*AssignOutput, apperrors.ApplicationError)
	}

	assignUsecase struct {
		contextFactory appcontext.Factory
	}

	AssignOutput struct {
		Data OperationResultData `json:"data"`
	}
)

func NewAssignUsecase(contextFactory appcontext.Factory) AssignUsecase {
	return &assignUsecase{contextFactory: contextFactory}
}

func (u *assignUsecase) Execute(ctx context.Context, teacherID, studentID string) (*AssignOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	if err := app.Repositories.TeacherStudentAssignment.Assign(ctx, domain.TeacherStudentAssignment{
		TeacherID: teacherID,
		StudentID: studentID,
		Status:    "active",
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentCreateError, err)
	}
	return &AssignOutput{Data: toOperationResultData(domain.OperationResult{Message: "student assigned successfully"})}, nil
}
