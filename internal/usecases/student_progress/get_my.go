package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GetMyProgressUsecase interface {
	Execute(context.Context, string) (*ProgressListOutput, apperrors.ApplicationError)
}

type getMyProgressUsecase struct {
	factory appcontext.Factory
}

func NewGetMyProgressUsecase(factory appcontext.Factory) GetMyProgressUsecase {
	return &getMyProgressUsecase{factory: factory}
}

func (u *getMyProgressUsecase) Execute(ctx context.Context, studentID string) (*ProgressListOutput, apperrors.ApplicationError) {
	app := u.factory()

	list, err := app.Repositories.StudentProgress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	var data []ProgressData
	for _, p := range list {
		data = append(data, toProgressData(p))
	}
	if data == nil {
		data = []ProgressData{}
	}

	return &ProgressListOutput{Data: data}, nil
}
