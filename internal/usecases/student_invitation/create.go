package studentinvitation

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/invitecode"
)

// defaultTTL: el código se dicta al principio del ciclo lectivo y tiene que
// durar el trimestre sin que el docente lo renueve.
const defaultTTL = 90 * 24 * time.Hour

// createAttempts cubre la colisión de códigos, que con 40 bits es rarísima
// pero no imposible; el índice único es el que la detecta.
const createAttempts = 3

type (
	CreateUsecase interface {
		Execute(ctx context.Context, teacherID string) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateOutput struct {
		Data InvitationData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

// Execute deja al docente con un único código vigente: regenerar revoca el
// anterior, así el que se filtró deja de servir apenas se pide otro.
func (u *createUsecase) Execute(ctx context.Context, teacherID string) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	current, err := app.Repositories.StudentInvitation.GetActiveByTeacher(ctx, teacherID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InvitationGetError, err)
	}
	if current != nil {
		if err := app.Repositories.StudentInvitation.Revoke(ctx, current.ID, teacherID); err != nil {
			return nil, apperrors.NewApplicationError(mappings.InvitationRevokeError, err)
		}
	}

	expiresAt := time.Now().Add(defaultTTL)

	var lastErr error
	for range createAttempts {
		code, err := invitecode.New()
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.InvitationCreateError, err)
		}

		created, err := app.Repositories.StudentInvitation.Create(ctx, domain.StudentInvitation{
			Code:      code,
			TeacherID: teacherID,
			ExpiresAt: &expiresAt,
		})
		if err == nil {
			return &CreateOutput{Data: toInvitationData(*created)}, nil
		}
		lastErr = err
	}

	return nil, apperrors.NewApplicationError(mappings.InvitationCreateError, lastErr)
}
