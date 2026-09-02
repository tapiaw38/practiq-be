package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	RemoveMemberUsecase interface {
		Execute(context.Context, string, string) (*RemoveMemberOutput, apperrors.ApplicationError)
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

func (u *removeMemberUsecase) Execute(ctx context.Context, gradeID, userID string) (*RemoveMemberOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	if err := app.Repositories.Grade.RemoveMember(ctx, gradeID, userID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeAssignMemberError, err)
	}
	return &RemoveMemberOutput{Data: toOperationResultData(domain.OperationResult{Message: "member removed successfully"})}, nil
}
