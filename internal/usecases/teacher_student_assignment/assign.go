package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
	"github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

type (
	AssignUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, in AssignInput) (*AssignOutput, apperrors.ApplicationError)
	}

	assignUsecase struct {
		contextFactory appcontext.Factory
	}

	AssignInput struct {
		TeacherID string `json:"teacher_id" binding:"required"`
		StudentID string `json:"student_id" binding:"required"`
	}

	AssignOutput struct {
		Data OperationResultData `json:"data"`
	}
)

func NewAssignUsecase(contextFactory appcontext.Factory) AssignUsecase {
	return &assignUsecase{contextFactory: contextFactory}
}

func (u *assignUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, in AssignInput) (*AssignOutput, apperrors.ApplicationError) {
	teacherID, studentID := in.TeacherID, in.StudentID
	app := u.contextFactory()

	if appErr := school.EnsureCanLinkTeacherStudent(ctx, app, requesterID, isSuperAdmin, teacherID, studentID); appErr != nil {
		return nil, appErr
	}

	if appErr := subscription.EnsureCanAddStudent(ctx, app, "", teacherID, studentID); appErr != nil {
		return nil, appErr
	}

	if err := app.Repositories.TeacherStudentAssignment.Assign(ctx, domain.TeacherStudentAssignment{
		TeacherID: teacherID,
		StudentID: studentID,
		Status:    "active",
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.AssignmentCreateError, err)
	}
	school.JoinTeacherSchool(ctx, app, teacherID, studentID)

	return &AssignOutput{Data: toOperationResultData(domain.OperationResult{Message: "student assigned successfully"})}, nil
}
