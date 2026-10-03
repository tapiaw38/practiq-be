package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	ListMembersUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID, bearerToken string) (*ListMembersOutput, apperrors.ApplicationError)
	}

	listMembersUsecase struct {
		contextFactory appcontext.Factory
	}

	ListMembersOutput struct {
		Data []MemberData `json:"data"`
	}
)

func NewListMembersUsecase(contextFactory appcontext.Factory) ListMembersUsecase {
	return &listMembersUsecase{contextFactory: contextFactory}
}

func (u *listMembersUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID, bearerToken string) (*ListMembersOutput, apperrors.ApplicationError) {
	data, appErr := listMembers(ctx, u.contextFactory(), requesterID, isSuperAdmin, schoolID, bearerToken)
	if appErr != nil {
		return nil, appErr
	}
	return &ListMembersOutput{Data: data}, nil
}

func listMembers(ctx context.Context, app *appcontext.Context, requesterID string, isSuperAdmin bool, schoolID, bearerToken string) ([]MemberData, apperrors.ApplicationError) {
	if !isSuperAdmin {
		if appErr := EnsureAdministers(ctx, app, requesterID, false, schoolID); appErr != nil {
			return nil, appErr
		}
	}
	members, err := app.Repositories.School.ListMembers(ctx, schoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	ids := make([]string, 0, len(members))
	for _, member := range members {
		ids = append(ids, member.UserID)
	}
	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ids)
	if appErr != nil {
		return nil, appErr
	}
	data := make([]MemberData, 0, len(members))
	for _, member := range members {
		info := names[member.UserID]
		data = append(data, MemberData{UserID: member.UserID, Name: identity.FullName(info, member.UserID), Email: info.Email, Role: member.Role, Active: member.Active})
	}
	return data, nil
}
