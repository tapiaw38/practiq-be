package school

import (
	"context"
	"regexp"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type Service interface {
	Create(context.Context, CreateInput) (*Output, apperrors.ApplicationError)
	List(context.Context, string, bool) (*ListOutput, apperrors.ApplicationError)
	Members(context.Context, string, string, bool) (*MembersOutput, apperrors.ApplicationError)
	AddMember(context.Context, AddMemberInput) (*MemberOutput, apperrors.ApplicationError)
	RemoveMember(context.Context, RemoveMemberInput) apperrors.ApplicationError
}

type service struct{ contextFactory appcontext.Factory }

type CreateInput struct{ Name, Slug, CreatedBy string }
type AddMemberInput struct {
	SchoolID, UserID, MembershipRole, ProfileType, ActorID string
	IsSuperAdmin                                           bool
}
type RemoveMemberInput struct {
	SchoolID, UserID, ActorID string
	IsSuperAdmin              bool
}
type Output struct {
	Data domain.School `json:"data"`
}
type ListOutput struct {
	Data []domain.School `json:"data"`
}
type MembersOutput struct {
	Data []domain.SchoolMembership `json:"data"`
}
type MemberOutput struct {
	Data domain.SchoolMembership `json:"data"`
}

func NewService(f appcontext.Factory) Service { return &service{contextFactory: f} }

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func (s *service) Create(ctx context.Context, in CreateInput) (*Output, apperrors.ApplicationError) {
	name := strings.TrimSpace(in.Name)
	slug := strings.ToLower(strings.TrimSpace(in.Slug))
	if name == "" || !slugPattern.MatchString(slug) {
		return nil, apperrors.NewBadRequestError("name and valid slug are required")
	}
	created, err := s.contextFactory().Repositories.School.Create(ctx, domain.School{Name: name, Slug: slug, CreatedBy: in.CreatedBy})
	if err != nil {
		return nil, apperrors.NewApplicationError(structuredError("school:create"), err)
	}
	return &Output{Data: *created}, nil
}

func (s *service) List(ctx context.Context, userID string, global bool) (*ListOutput, apperrors.ApplicationError) {
	items, err := s.contextFactory().Repositories.School.List(ctx, userID, global)
	if err != nil {
		return nil, apperrors.NewApplicationError(structuredError("school:list"), err)
	}
	return &ListOutput{Data: items}, nil
}

func (s *service) Members(ctx context.Context, id, actor string, global bool) (*MembersOutput, apperrors.ApplicationError) {
	app := s.contextFactory()
	if _, err := app.Repositories.School.Get(ctx, id); err != nil {
		return nil, apperrors.NewApplicationError(structuredError("school:get"), err)
	}
	if !global {
		ok, err := app.Repositories.School.HasAdminAccess(ctx, id, actor)
		if err != nil {
			return nil, apperrors.NewApplicationError(structuredError("school:access"), err)
		}
		if !ok {
			return nil, apperrors.NewForbiddenError()
		}
	}
	items, err := app.Repositories.School.ListMembers(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(structuredError("school:members"), err)
	}
	return &MembersOutput{Data: items}, nil
}

func (s *service) AddMember(ctx context.Context, in AddMemberInput) (*MemberOutput, apperrors.ApplicationError) {
	app := s.contextFactory()
	if _, err := app.Repositories.School.Get(ctx, in.SchoolID); err != nil {
		return nil, apperrors.NewApplicationError(structuredError("school:get"), err)
	}
	if !in.IsSuperAdmin {
		ok, err := app.Repositories.School.HasAdminAccess(ctx, in.SchoolID, in.ActorID)
		if err != nil {
			return nil, apperrors.NewApplicationError(structuredError("school:access"), err)
		}
		if !ok {
			return nil, apperrors.NewForbiddenError()
		}
		if in.MembershipRole != "member" {
			return nil, apperrors.NewForbiddenError()
		}
	}
	if in.UserID == "" || (in.MembershipRole != "admin" && in.MembershipRole != "member") || (in.ProfileType != "student" && in.ProfileType != "teacher") {
		return nil, apperrors.NewBadRequestError("user_id, membership_role and profile_type are required")
	}
	m := domain.SchoolMembership{SchoolID: in.SchoolID, UserID: in.UserID, MembershipRole: in.MembershipRole, ProfileType: in.ProfileType}
	if err := app.Repositories.School.UpsertMember(ctx, m); err != nil {
		return nil, apperrors.NewApplicationError(structuredError("school:member"), err)
	}
	return &MemberOutput{Data: m}, nil
}

func (s *service) RemoveMember(ctx context.Context, in RemoveMemberInput) apperrors.ApplicationError {
	app := s.contextFactory()
	if !in.IsSuperAdmin {
		ok, err := app.Repositories.School.HasAdminAccess(ctx, in.SchoolID, in.ActorID)
		if err != nil {
			return apperrors.NewApplicationError(structuredError("school:access"), err)
		}
		if !ok {
			return apperrors.NewForbiddenError()
		}
	}
	if err := app.Repositories.School.DeleteMember(ctx, in.SchoolID, in.UserID); err != nil {
		return apperrors.NewApplicationError(structuredError("school:member"), err)
	}
	return nil
}

// Keep school errors local until the shared mappings gain a school namespace.
func structuredError(code string) (result mappings.ErrorDetails) {
	result.InternalCode, result.StatusCode, result.Message = code, 500, "school operation failed"
	return
}
