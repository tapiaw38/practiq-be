package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListMembersUsecase interface {
	Execute(context.Context, string) (*GradeMembersOutput, apperrors.ApplicationError)
}

type listMembersUsecase struct {
	factory appcontext.Factory
}

func NewListMembersUsecase(factory appcontext.Factory) ListMembersUsecase {
	return &listMembersUsecase{factory: factory}
}

func (u *listMembersUsecase) Execute(ctx context.Context, gradeID string) (*GradeMembersOutput, apperrors.ApplicationError) {
	app := u.factory()

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

	data := make([]GradeMemberData, 0, len(members))
	for _, member := range members {
		data = append(data, toGradeMemberData(member))
	}

	return &GradeMembersOutput{Data: data}, nil
}
