package userprofile

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	FindByEmailUsecase interface {
		Execute(ctx context.Context, email, bearerToken string) (*FindByEmailOutput, apperrors.ApplicationError)
	}

	findByEmailUsecase struct {
		contextFactory appcontext.Factory
	}

	FindByEmailOutput struct {
		Data FindByEmailData `json:"data"`
	}

	FindByEmailData struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

func NewFindByEmailUsecase(contextFactory appcontext.Factory) FindByEmailUsecase {
	return &findByEmailUsecase{contextFactory: contextFactory}
}

func (u *findByEmailUsecase) Execute(ctx context.Context, email, bearerToken string) (*FindByEmailOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, apperrors.NewBadRequestError("email is required")
	}

	info, appErr := identity.ByEmail(ctx, app.Integrations.AuthAPI, bearerToken, email)
	if appErr != nil {
		return nil, appErr
	}
	if info == nil {
		return nil, apperrors.NewNotFoundError("no Practiq account found for that email")
	}

	return &FindByEmailOutput{Data: FindByEmailData{
		ID:        info.Username,
		FirstName: info.FirstName,
		LastName:  info.LastName,
		Email:     info.Email,
	}}, nil
}
