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
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Name    string `json:"name"`
		Kind    string `json:"kind"`
		Billing string `json:"billing"`
	}

	UpdateOutput struct {
		Data SchoolData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	name, kind, billing := in.Name, in.Kind, in.Billing
	app := u.contextFactory()
	if appErr := EnsureAdministers(ctx, app, requesterID, isSuperAdmin, id); appErr != nil {
		return nil, appErr
	}
	current, err := app.Repositories.School.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if current == nil {
		return nil, apperrors.NewNotFoundError("school not found")
	}
	update := domain.School{Name: strings.TrimSpace(name)}
	if isSuperAdmin && current.Kind == domain.SchoolKindInstitution {
		update.Kind, update.Billing = kind, billing
		if update.Kind != "" && update.Kind != domain.SchoolKindInstitution {
			return nil, apperrors.NewBadRequestError("an institution cannot become a personal school")
		}
		if update.Billing != "" && update.Billing != domain.SchoolBillingDirect && update.Billing != domain.SchoolBillingSubscription {
			return nil, apperrors.NewBadRequestError("billing must be direct or subscription")
		}
	}
	if err := app.Repositories.School.Update(ctx, id, update); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	data, appErr := readSchool(ctx, app, id, "")
	if appErr != nil {
		return nil, appErr
	}
	return &UpdateOutput{Data: data}, nil
}
