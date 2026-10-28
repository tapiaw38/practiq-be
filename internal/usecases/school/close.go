package school

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	CloseUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in CloseInput) (*CloseOutput, apperrors.ApplicationError)
	}

	closeUsecase struct {
		contextFactory appcontext.Factory
	}

	CloseInput struct {
		ConfirmName string `json:"confirm_name"`
		Reason      string `json:"reason"`
	}

	CloseOutput struct {
		Data SchoolData `json:"data"`
	}
)

func NewCloseUsecase(contextFactory appcontext.Factory) CloseUsecase {
	return &closeUsecase{contextFactory: contextFactory}
}

func (u *closeUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in CloseInput) (*CloseOutput, apperrors.ApplicationError) {
	confirmName, reason := in.ConfirmName, in.Reason
	if !isSuperAdmin {
		return nil, apperrors.NewForbiddenError()
	}
	app := u.contextFactory()
	school, appErr := schoolForLifecycle(ctx, app, id)
	if appErr != nil {
		return nil, appErr
	}
	if school.Status == domain.SchoolStatusClosed {
		data, appErr := readSchool(ctx, app, id, "")
		if appErr != nil {
			return nil, appErr
		}
		return &CloseOutput{Data: data}, nil
	}
	if strings.TrimSpace(confirmName) != school.Name {
		return nil, apperrors.NewBadRequestError("confirmation must match the school name")
	}
	if err := app.Repositories.School.Close(ctx, id, requesterID, strings.TrimSpace(reason)); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if err := app.Repositories.StudentInvitation.RevokeForSchool(ctx, id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.InvitationRevokeError, err)
	}
	data, appErr := readSchool(ctx, app, id, "")
	if appErr != nil {
		return nil, appErr
	}
	return &CloseOutput{Data: data}, nil
}
