package userprofile

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
)

type (
	SyncUsecase interface {
		Execute(ctx context.Context, id, bearerToken string, in SyncInput) (*SyncOutput, apperrors.ApplicationError)
	}

	syncUsecase struct {
		contextFactory appcontext.Factory
	}

	SyncInput struct {
		ProfileType string `json:"profile_type"`
		Timezone    string `json:"timezone"`
	}

	SyncOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewSyncUsecase(contextFactory appcontext.Factory) SyncUsecase {
	return &syncUsecase{contextFactory: contextFactory}
}

func (u *syncUsecase) Execute(ctx context.Context, id, bearerToken string, in SyncInput) (*SyncOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	existing, err := app.Repositories.UserProfile.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	profileType := in.ProfileType
	if existing != nil {
		profileType = existing.ProfileType
	} else if profileType == "" {
		profileType = "student"
	}
	if profileType != "teacher" && profileType != "student" {
		return nil, apperrors.NewBadRequestError("profile_type must be teacher or student")
	}

	p := domain.UserProfile{
		ID:          id,
		ProfileType: profileType,
		Timezone:    in.Timezone,
	}

	if err := app.Repositories.UserProfile.Upsert(ctx, p); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}

	updated, err := app.Repositories.UserProfile.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{id})
	if appErr != nil {
		return nil, appErr
	}
	info := names[id]
	displayName := identity.FullName(info, id)

	ensurePersonalSchool(ctx, app, *updated, displayName)

	return &SyncOutput{Data: toProfileData(*updated, displayName, info.Email, assistantcfg.Enabled(ctx, app))}, nil
}

func ensurePersonalSchool(ctx context.Context, app *appcontext.Context, profile domain.UserProfile, displayName string) {
	if profile.ProfileType != "teacher" {
		return
	}

	existing, err := app.Repositories.School.GetPersonal(ctx, profile.ID)
	if err != nil {
		log.Printf("[schools] personal school lookup failed user_id=%s err=%v", profile.ID, err)
		return
	}

	wanted := domain.PersonalSchoolName(displayName)

	if existing != nil {

		if existing.Name == domain.PlaceholderSchoolName && wanted != domain.PlaceholderSchoolName {
			if err := app.Repositories.School.Rename(ctx, existing.ID, wanted); err != nil {
				log.Printf("[schools] rename failed school_id=%s err=%v", existing.ID, err)
			}
		}
		return
	}

	schoolID, err := app.Repositories.School.Create(ctx, domain.School{
		Name:      wanted,
		Kind:      domain.SchoolKindPersonal,
		Billing:   domain.SchoolBillingSubscription,
		CreatedBy: profile.ID,
	})
	if err != nil {
		log.Printf("[schools] create failed user_id=%s err=%v", profile.ID, err)
		return
	}

	if err := app.Repositories.School.AddMember(ctx, domain.SchoolMember{
		SchoolID: schoolID,
		UserID:   profile.ID,
		Role:     domain.SchoolRoleAdmin,
	}); err != nil {
		log.Printf("[schools] membership failed school_id=%s err=%v", schoolID, err)
	}
}
