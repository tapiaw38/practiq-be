package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListUsecase interface {
	Execute(context.Context) (*SubjectListOutput, apperrors.ApplicationError)
}

type listUsecase struct {
	factory appcontext.Factory
}

func NewListUsecase(factory appcontext.Factory) ListUsecase {
	return &listUsecase{factory: factory}
}

func (u *listUsecase) Execute(ctx context.Context) (*SubjectListOutput, apperrors.ApplicationError) {
	app := u.factory()
	subjects, err := app.Repositories.Subject.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectListError, err)
	}
	data := make([]SubjectData, 0, len(subjects))
	for _, subject := range subjects {
		data = append(data, toSubjectData(subject))
	}
	return &SubjectListOutput{Data: data}, nil
}
