package studentinvitation

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	RevokeUsecase interface {
		Execute(ctx context.Context, id, teacherID string) apperrors.ApplicationError
	}

	revokeUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewRevokeUsecase(contextFactory appcontext.Factory) RevokeUsecase {
	return &revokeUsecase{contextFactory: contextFactory}
}

func (u *revokeUsecase) Execute(ctx context.Context, id, teacherID string) apperrors.ApplicationError {
	app := u.contextFactory()

	if err := app.Repositories.StudentInvitation.Revoke(ctx, id, teacherID); err != nil {
		return apperrors.NewApplicationError(mappings.InvitationRevokeError, err)
	}

	return nil
}
