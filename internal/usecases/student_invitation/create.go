package studentinvitation

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/invitecode"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

const defaultTTL = 90 * 24 * time.Hour

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

	schoolID, appErr := school.OwnedSchoolID(ctx, app, teacherID)
	if appErr != nil {
		return nil, appErr
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
			SchoolID:  schoolID,
			ExpiresAt: &expiresAt,
		})
		if err == nil {
			return &CreateOutput{Data: toInvitationData(*created)}, nil
		}
		lastErr = err
	}

	return nil, apperrors.NewApplicationError(mappings.InvitationCreateError, lastErr)
}
