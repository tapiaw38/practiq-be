package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	RemoveMemberUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID, userID string) apperrors.ApplicationError
	}

	removeMemberUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewRemoveMemberUsecase(contextFactory appcontext.Factory) RemoveMemberUsecase {
	return &removeMemberUsecase{contextFactory: contextFactory}
}

func (u *removeMemberUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID, userID string) apperrors.ApplicationError {
	app := u.contextFactory()
	if appErr := EnsureAdministers(ctx, app, requesterID, isSuperAdmin, schoolID); appErr != nil {
		return appErr
	}
	members, err := app.Repositories.School.ListMembers(ctx, schoolID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	for _, member := range members {
		if member.UserID == userID && member.Role == domain.SchoolRoleAdmin && member.Active {
			admins, countErr := app.Repositories.School.CountActiveAdmins(ctx, schoolID)
			if countErr != nil {
				return apperrors.NewApplicationError(mappings.SchoolLookupError, countErr)
			}
			if admins <= 1 {
				return apperrors.NewBadRequestError("an active school needs at least one admin")
			}
			break
		}
	}
	if !isSuperAdmin && userID == requesterID {
		return apperrors.NewBadRequestError("you cannot remove yourself from a school you administer")
	}
	if err := app.Repositories.School.RemoveMember(ctx, schoolID, userID); err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return nil
}
