package sitecontact

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
		Execute(ctx context.Context, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone" binding:"required"`
		WhatsApp string `json:"whatsapp" binding:"required"`
	}

	UpdateOutput struct {
		Data ContactData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	saved, err := app.Repositories.SiteContact.Save(ctx, domain.SiteContact{
		Email:    strings.TrimSpace(in.Email),
		Phone:    strings.TrimSpace(in.Phone),
		WhatsApp: strings.TrimSpace(in.WhatsApp),
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SiteContactSaveError, err)
	}
	return &UpdateOutput{Data: toContactData(saved)}, nil
}
