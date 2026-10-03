package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	MineUsecase interface {
		Execute(ctx context.Context, requesterID, bearerToken string) (*MineOutput, apperrors.ApplicationError)
	}

	mineUsecase struct {
		contextFactory appcontext.Factory
	}

	MineOutput struct {
		Data []SchoolData `json:"data"`
	}
)

func NewMineUsecase(contextFactory appcontext.Factory) MineUsecase {
	return &mineUsecase{contextFactory: contextFactory}
}

func (u *mineUsecase) Execute(ctx context.Context, requesterID, bearerToken string) (*MineOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	members, err := app.Repositories.School.ListForUser(ctx, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	schools := make([]domain.School, 0, len(members))
	roles := make(map[string]string, len(members))
	for _, member := range members {
		if !member.Active {
			continue
		}
		school, err := app.Repositories.School.Get(ctx, member.SchoolID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
		}
		if school != nil {
			schools, roles[school.ID] = append(schools, *school), member.Role
		}
	}
	schools, appErr := resolvePersonalSchoolNames(ctx, app, bearerToken, schools)
	if appErr != nil {
		return nil, appErr
	}
	data := make([]SchoolData, 0, len(schools))
	for _, school := range schools {
		data = append(data, toSchoolData(school, roles[school.ID]))
	}
	return &MineOutput{Data: data}, nil
}
