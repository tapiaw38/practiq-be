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
	CreateUsecase interface {
		Execute(ctx context.Context, in CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Name        string `json:"name"`
		Kind        string `json:"kind"`
		Billing     string `json:"billing"`
		AdminUserID string `json:"admin_user_id"`
	}

	CreateOutput struct {
		Data SchoolData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, in CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	name, kind, billing, adminUserID := in.Name, in.Kind, in.Billing, in.AdminUserID
	app := u.contextFactory()
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperrors.NewBadRequestError("a school needs a name")
	}
	if kind == "" {
		kind = domain.SchoolKindInstitution
	}
	if billing == "" {
		billing = domain.SchoolBillingDirect
	}
	if kind != domain.SchoolKindInstitution {
		return nil, apperrors.NewBadRequestError("only institutions can be created here")
	}
	if billing != domain.SchoolBillingDirect && billing != domain.SchoolBillingSubscription {
		return nil, apperrors.NewBadRequestError("billing must be direct or subscription")
	}
	adminID := strings.TrimSpace(adminUserID)
	if adminID == "" {
		return nil, apperrors.NewBadRequestError("an institution needs its first administrator")
	}
	profile, err := app.Repositories.UserProfile.Get(ctx, adminID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil {
		return nil, apperrors.NewNotFoundError("the administrator has no Practiq profile yet")
	}
	if profile.ProfileType != "teacher" {
		return nil, apperrors.NewBadRequestError("the first administrator must have a teacher profile")
	}
	id, err := app.Repositories.School.CreateWithAdmin(ctx, domain.School{Name: name, Kind: kind, Billing: billing}, adminID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	data, appErr := readSchool(ctx, app, id, "")
	if appErr != nil {
		return nil, appErr
	}
	return &CreateOutput{Data: data}, nil
}
