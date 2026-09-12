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
		Execute(context.Context, string, bool, string, string) (*AssignOutput, apperrors.ApplicationError)
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

func (u *assignUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, teacherID, studentID string) (*AssignOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	// A school admin may only link people who already share their school —
	// see school.EnsureCanLinkTeacherStudent for why this used to be
	// superadmin-only.
	if appErr := school.EnsureCanLinkTeacherStudent(ctx, app, requesterID, isSuperAdmin, teacherID, studentID); appErr != nil {
		return nil, appErr
	}
	// A superadmin assigning is still a teacher gaining a student, and the
	// plan is the teacher's either way. No school is named here — a direct
	// assignment is the teacher's own — so the cap comes from theirs.
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
