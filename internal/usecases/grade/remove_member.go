package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type RemoveMemberUsecase interface {
	Execute(context.Context, string, string) (*RemoveMemberOutput, apperrors.ApplicationError)
}

type removeMemberUsecase struct {
	factory appcontext.Factory
}

func NewRemoveMemberUsecase(factory appcontext.Factory) RemoveMemberUsecase {
	return &removeMemberUsecase{factory: factory}
}

func (u *removeMemberUsecase) Execute(ctx context.Context, gradeID, userID string) (*RemoveMemberOutput, apperrors.ApplicationError) {
	app := u.factory()
	if err := app.Repositories.Grade.RemoveMember(ctx, gradeID, userID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeAssignMemberError, err)
	}
	return &RemoveMemberOutput{Message: "member removed successfully"}, nil
}
