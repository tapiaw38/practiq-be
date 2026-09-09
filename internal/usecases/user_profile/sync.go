package userprofile

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	SyncUsecase interface {
		Execute(context.Context, SyncInput) (*SyncOutput, apperrors.ApplicationError)
	}

	syncUsecase struct {
		contextFactory appcontext.Factory
	}

	SyncInput struct {
		ID               string
		ProfileType      string
		Timezone         string
		AssistantBaseURL string
		AssistantAPIKey  string
		// BearerToken is the caller's own "Bearer <jwt>" header, forwarded to
		// auth-api-be to resolve the caller's own display name — never
		// trusted from the request body, since that would let the client
		// claim any name it likes.
		BearerToken string
	}

	SyncOutput struct {
		Data ProfileData `json:"data"`
	}
)

func NewSyncUsecase(contextFactory appcontext.Factory) SyncUsecase {
	return &syncUsecase{contextFactory: contextFactory}
}

func (u *syncUsecase) Execute(ctx context.Context, input SyncInput) (*SyncOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	existing, err := app.Repositories.UserProfile.Get(ctx, input.ID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	// A profile type is selected exactly once during Practiq onboarding. This
	// stops a student from turning into a teacher by replaying POST /profile,
	// while keeping Auth roles out of product authorization.
	profileType := input.ProfileType
	if existing != nil {
		profileType = existing.ProfileType
	} else if profileType == "" {
		profileType = "student"
	}
	if profileType != "teacher" && profileType != "student" {
		return nil, apperrors.NewBadRequestError("profile_type must be teacher or student")
	}

	p := domain.UserProfile{
		ID:               input.ID,
		ProfileType:      profileType,
		Timezone:         input.Timezone,
		AssistantBaseURL: input.AssistantBaseURL,
		AssistantAPIKey:  input.AssistantAPIKey,
	}

	if err := app.Repositories.UserProfile.Upsert(ctx, p); err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileSyncError, err)
	}

	updated, err := app.Repositories.UserProfile.Get(ctx, input.ID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	names, err := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, []string{input.ID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	info := names[input.ID]
	displayName := identity.FullName(info, input.ID)

	ensurePersonalSchool(ctx, app, *updated, displayName)

	return &SyncOutput{Data: toProfileData(*updated, displayName, info.Email)}, nil
}

// ensurePersonalSchool gives a teacher the school they administer.
//
// It runs on every sync rather than only on the first, because the schools
// migration could not name these: teacher names live in auth-api-be, so the SQL
// left a placeholder. The first sync that can resolve a real name replaces it.
//
// Failures are logged and swallowed. A teacher who cannot log in because their
// school could not be created is a worse outcome than one who logs in and gets
// it on the next sync, and this runs on every request that touches the profile.
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
		// Only the placeholder is replaced. A teacher who renamed their school
		// keeps that name; overwriting it on every login would undo their edit.
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
