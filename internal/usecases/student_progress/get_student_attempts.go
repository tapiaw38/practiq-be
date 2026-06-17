package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetStudentAttemptsUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, studentID, sheetID string) (*GetStudentAttemptsOutput, apperrors.ApplicationError)
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

func (u *getStudentAttemptsUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, studentID, sheetID string) (*GetStudentAttemptsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Authorization check: only admins or teachers with access can view student attempts
	if !isAdmin {
		hasAccess, err := app.Repositories.TeacherStudentAssignment.HasAccess(ctx, requesterID, studentID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.AssignmentListError, err)
		}
		if !hasAccess {
			return nil, apperrors.NewForbiddenError()
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
