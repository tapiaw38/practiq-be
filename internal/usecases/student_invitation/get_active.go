package studentinvitation

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetActiveUsecase interface {
		Execute(ctx context.Context, teacherID string) (*GetActiveOutput, apperrors.ApplicationError)
	}

	getActiveUsecase struct {
		contextFactory appcontext.Factory
	}

	// Data en nil significa que el docente todavía no generó ninguno, o que el
	// que tenía venció. La pantalla ofrece generarlo en los dos casos.
	GetActiveOutput struct {
		Data *InvitationData `json:"data"`
	}
)

func NewGetActiveUsecase(contextFactory appcontext.Factory) GetActiveUsecase {
	return &getActiveUsecase{contextFactory: contextFactory}
}

func (u *getActiveUsecase) Execute(ctx context.Context, teacherID string) (*GetActiveOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	invitation, err := app.Repositories.StudentInvitation.GetActiveByTeacher(ctx, teacherID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InvitationGetError, err)
	}
	if invitation == nil || !invitation.IsUsable(time.Now()) {
		return &GetActiveOutput{}, nil
	}

	data := toInvitationData(*invitation)

	return &GetActiveOutput{Data: &data}, nil
}
