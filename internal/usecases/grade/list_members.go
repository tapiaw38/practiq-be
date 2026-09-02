package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	ListMembersUsecase interface {
		Execute(ctx context.Context, gradeID, bearerToken string) (*ListMembersOutput, apperrors.ApplicationError)
	}

	listMembersUsecase struct {
		contextFactory appcontext.Factory
	}

	ListMembersOutput struct {
		Data []GradeMemberData `json:"data"`
	}
)

func NewListMembersUsecase(contextFactory appcontext.Factory) ListMembersUsecase {
	return &listMembersUsecase{contextFactory: contextFactory}
}

func (u *listMembersUsecase) Execute(ctx context.Context, gradeID, bearerToken string) (*ListMembersOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	grade, err := app.Repositories.Grade.Get(ctx, gradeID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if grade == nil {
		return nil, apperrors.NewApplicationError(mappings.GradeNotFoundError, nil)
	}

	members, err := app.Repositories.Grade.ListMembers(ctx, gradeID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeListMembersError, err)
	}

	ids := make([]string, 0, len(members))
	for _, member := range members {
		ids = append(ids, member.ID)
	}
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ids)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	data := make([]GradeMemberData, 0, len(members))
	for _, member := range members {
		info := names[member.ID]
		data = append(data, toGradeMemberData(member, identity.FullName(info, member.ID), info.Email))
	}

	return &ListMembersOutput{Data: data}, nil
}
