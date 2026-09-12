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

// Execute no distingue entre una invitación ajena y una inexistente: el
// repositorio filtra por teacher_id y en ambos casos no cambia ninguna fila.
// Los alumnos ya vinculados siguen vinculados; revocar corta el ingreso de
// nuevos, para desvincular está la baja de la asignación.
func (u *revokeUsecase) Execute(ctx context.Context, id, teacherID string) apperrors.ApplicationError {
	app := u.contextFactory()

	if err := app.Repositories.StudentInvitation.Revoke(ctx, id, teacherID); err != nil {
		return apperrors.NewApplicationError(mappings.InvitationRevokeError, err)
	}

	return nil
}
