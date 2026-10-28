package sitecontact

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(ctx context.Context) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data ContactData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	contact, err := app.Repositories.SiteContact.Get(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SiteContactGetError, err)
	}
	return &GetOutput{Data: toContactData(contact)}, nil
}
