package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	ListUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID string) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []SubjectData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	scope, err := school.ScopeForSchool(ctx, app, requesterID, isSuperAdmin, schoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectListError, err)
	}
	if scope.Empty() {
		// Belonging to no school means seeing nothing, not everything.
		return &ListOutput{Data: []SubjectData{}}, nil
	}
	subjects, err := app.Repositories.Subject.List(ctx, scope.SchoolIDs)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectListError, err)
	}
	data := make([]SubjectData, 0, len(subjects))
	for _, subject := range subjects {
		data = append(data, toSubjectData(subject))
	}
	return &ListOutput{Data: data}, nil
}
