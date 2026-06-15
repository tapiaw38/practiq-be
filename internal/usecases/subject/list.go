package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context) (*ListOutput, apperrors.ApplicationError)
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

func (u *listUsecase) Execute(ctx context.Context) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	subjects, err := app.Repositories.Subject.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectListError, err)
	}
	data := make([]SubjectData, 0, len(subjects))
	for _, subject := range subjects {
		data = append(data, toSubjectData(subject))
	}
	return &ListOutput{Data: data}, nil
}
