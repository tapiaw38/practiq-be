package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	schoolUC "github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	RemoveMemberUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, gradeID, userID string) (*RemoveMemberOutput, apperrors.ApplicationError)
	}

	removeMemberUsecase struct {
		contextFactory appcontext.Factory
	}

	RemoveMemberOutput struct {
		Data OperationResultData `json:"data"`
	}
)

func NewRemoveMemberUsecase(contextFactory appcontext.Factory) RemoveMemberUsecase {
	return &removeMemberUsecase{contextFactory: contextFactory}
}

func (u *removeMemberUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, gradeID, userID string) (*RemoveMemberOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// The route lets a teacher ask; this decides whose grade they may touch.
	grade, err := app.Repositories.Grade.Get(ctx, gradeID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if grade == nil {
		return nil, apperrors.NewNotFoundError("grade not found")
	}
	if appErr := schoolUC.EnsureAdministers(ctx, app, requesterID, isSuperAdmin, grade.SchoolID); appErr != nil {
		return nil, appErr
	}

	if err := app.Repositories.Grade.RemoveMember(ctx, gradeID, userID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeAssignMemberError, err)
	}
	return &RemoveMemberOutput{Data: toOperationResultData(domain.OperationResult{Message: "member removed successfully"})}, nil
}
