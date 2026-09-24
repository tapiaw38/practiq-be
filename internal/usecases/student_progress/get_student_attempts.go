package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	GetStudentAttemptsUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, studentID, sheetID string) (*GetStudentAttemptsOutput, apperrors.ApplicationError)
	}

	getStudentAttemptsUsecase struct {
		contextFactory appcontext.Factory
	}

	GetStudentAttemptsOutput struct {
		Data []AttemptData `json:"data"`
	}
)

func NewGetStudentAttemptsUsecase(contextFactory appcontext.Factory) GetStudentAttemptsUsecase {
	return &getStudentAttemptsUsecase{contextFactory: contextFactory}
}

func (u *getStudentAttemptsUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, studentID, sheetID string) (*GetStudentAttemptsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if !isSuperAdmin {
		hasAccess, err := app.Repositories.TeacherStudentAssignment.HasAccess(ctx, requesterID, studentID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
		}
		if !hasAccess {
			return nil, apperrors.NewForbiddenError()
		}
	}

	// HasAccess only proves the requester can see this student somewhere; it
	// says nothing about the sheet. Without tying the sheet back to a course
	// the requester owns, a teacher of one course could read the same student's
	// answers and scores for a sheet in another course, or in a deleted one.
	if !isSuperAdmin {
		sheet, err := app.Repositories.PracticeSheet.Get(ctx, sheetID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
		}
		if sheet == nil {
			return nil, apperrors.NewNotFoundError("practice sheet not found")
		}
		// EnsureCanManageCourse's Get filters out soft-deleted courses, so a
		// not-found error here is also the deleted-course case.
		if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, sheet.CourseID); appErr != nil {
			return nil, appErr
		}
	}

	attempts, err := app.Repositories.StudentAttempt.ListBySheet(ctx, studentID, sheetID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	var data []AttemptData
	for _, a := range attempts {
		data = append(data, toAttemptData(a))
	}
	if data == nil {
		data = []AttemptData{}
	}

	return &GetStudentAttemptsOutput{Data: data}, nil
}
