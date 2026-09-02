package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UnassignUsecase interface {
		Execute(context.Context, string, string) (*UnassignOutput, apperrors.ApplicationError)
	}

	unassignUsecase struct {
		contextFactory appcontext.Factory
	}

	UnassignOutput struct {
		Data OperationResultData `json:"data"`
	}
)

func NewUnassignUsecase(contextFactory appcontext.Factory) UnassignUsecase {
	return &unassignUsecase{contextFactory: contextFactory}
}

func (u *unassignUsecase) Execute(ctx context.Context, teacherID, studentID string) (*UnassignOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	if err := app.Repositories.TeacherStudentAssignment.Unassign(ctx, teacherID, studentID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentDeleteError, err)
	}
	return &UnassignOutput{Data: toOperationResultData(domain.OperationResult{Message: "student unassigned successfully"})}, nil
}
