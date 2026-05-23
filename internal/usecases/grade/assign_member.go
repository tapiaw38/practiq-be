package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type AssignMemberUsecase interface {
	Execute(context.Context, string, string) (*AssignMemberOutput, apperrors.ApplicationError)
}

type assignMemberUsecase struct {
	factory appcontext.Factory
}

func NewAssignMemberUsecase(factory appcontext.Factory) AssignMemberUsecase {
	return &assignMemberUsecase{factory: factory}
}

func (u *assignMemberUsecase) Execute(ctx context.Context, gradeID, userID string) (*AssignMemberOutput, apperrors.ApplicationError) {
	app := u.factory()

	grade, err := app.Repositories.Grade.Get(ctx, gradeID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if grade == nil {
		return nil, apperrors.NewApplicationError(mappings.GradeNotFoundError, nil)
	}

	profile, err := app.Repositories.UserProfile.Get(ctx, userID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil {
		return nil, apperrors.NewApplicationError(mappings.NotFoundError, nil)
	}

	if err := app.Repositories.Grade.AddMember(ctx, domain.GradeMembership{
		GradeID: gradeID,
		UserID:  userID,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeAssignMemberError, err)
	}

	return &AssignMemberOutput{Message: "member assigned successfully"}, nil
}
