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
	AssignMemberUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, gradeID, userID string) (*AssignMemberOutput, apperrors.ApplicationError)
	}

	assignMemberUsecase struct {
		contextFactory appcontext.Factory
	}

	AssignMemberOutput struct {
		Data OperationResultData `json:"data"`
	}
)

func NewAssignMemberUsecase(contextFactory appcontext.Factory) AssignMemberUsecase {
	return &assignMemberUsecase{contextFactory: contextFactory}
}

func (u *assignMemberUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, gradeID, userID string) (*AssignMemberOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	grade, err := app.Repositories.Grade.Get(ctx, gradeID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if grade == nil {
		return nil, apperrors.NewApplicationError(mappings.GradeNotFoundError, nil)
	}

	// The route lets a teacher ask; this decides whose grade they may touch.
	if appErr := schoolUC.EnsureAdministers(ctx, app, requesterID, isSuperAdmin, grade.SchoolID); appErr != nil {
		return nil, appErr
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

	return &AssignMemberOutput{Data: toOperationResultData(domain.OperationResult{Message: "member assigned successfully"})}, nil
}
