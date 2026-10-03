package school

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

type (
	AddMemberUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID string, in AddMemberInput) apperrors.ApplicationError
	}

	addMemberUsecase struct {
		contextFactory appcontext.Factory
	}

	AddMemberInput struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
)

func NewAddMemberUsecase(contextFactory appcontext.Factory) AddMemberUsecase {
	return &addMemberUsecase{contextFactory: contextFactory}
}

func (u *addMemberUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID string, in AddMemberInput) apperrors.ApplicationError {
	userID, role := in.UserID, in.Role
	app := u.contextFactory()
	if appErr := EnsureAdministers(ctx, app, requesterID, isSuperAdmin, schoolID); appErr != nil {
		return appErr
	}
	if strings.TrimSpace(userID) == "" {
		return apperrors.NewBadRequestError("a member needs a user")
	}
	switch role {
	case domain.SchoolRoleAdmin, domain.SchoolRoleTeacher, domain.SchoolRoleStudent:
	default:
		return apperrors.NewBadRequestError("role must be admin, teacher or student")
	}
	school, err := app.Repositories.School.Get(ctx, schoolID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil {
		return apperrors.NewNotFoundError("school not found")
	}
	if school.Kind == domain.SchoolKindPersonal && role != domain.SchoolRoleStudent {
		return apperrors.NewForbiddenError()
	}
	members, err := app.Repositories.School.ListMembers(ctx, schoolID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	wasActiveStudent := false
	for _, member := range members {
		if member.UserID != userID {
			continue
		}
		wasActiveStudent = member.Role == domain.SchoolRoleStudent && member.Active
		if member.Role == domain.SchoolRoleAdmin && member.Active && role != domain.SchoolRoleAdmin {
			admins, countErr := app.Repositories.School.CountActiveAdmins(ctx, schoolID)
			if countErr != nil {
				return apperrors.NewApplicationError(mappings.SchoolLookupError, countErr)
			}
			if admins <= 1 {
				return apperrors.NewBadRequestError("an active school needs at least one admin")
			}
		}
		break
	}
	if role == domain.SchoolRoleStudent && !wasActiveStudent {
		if appErr := subscription.EnsureCanAddStudent(ctx, app, schoolID, requesterID, userID); appErr != nil {
			return appErr
		}
	}
	profile, err := app.Repositories.UserProfile.Get(ctx, userID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil {
		return apperrors.NewNotFoundError("that person has no Practiq profile yet — they have to sign in once before joining a school")
	}
	if role != domain.SchoolRoleStudent && profile.ProfileType != "teacher" {
		return apperrors.NewBadRequestError("an administrator or teacher must have a teacher profile")
	}
	if err := app.Repositories.School.AddMember(ctx, domain.SchoolMember{SchoolID: schoolID, UserID: userID, Role: role}); err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return nil
}
