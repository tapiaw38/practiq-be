package grade

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	ListUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []GradeData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	scope, err := school.ScopeFor(ctx, app, requesterID, isSuperAdmin)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeListError, err)
	}
	if scope.Empty() {
		// Belonging to no school means seeing nothing, not everything.
		return &ListOutput{Data: []GradeData{}}, nil
	}

	grades, err := app.Repositories.Grade.List(ctx, scope.SchoolIDs)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeListError, err)
	}

	data := make([]GradeData, 0, len(grades))
	for _, grade := range grades {
		data = append(data, toGradeData(grade))
	}

	return &ListOutput{Data: data}, nil
}
