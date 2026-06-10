package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type ListByCourseUsecase interface {
	Execute(context.Context, string) (*CourseLearningStrategyListOutput, apperrors.ApplicationError)
}

type listByCourseUsecase struct {
	factory appcontext.Factory
}

func NewListByCourseUsecase(factory appcontext.Factory) ListByCourseUsecase {
	return &listByCourseUsecase{factory: factory}
}

func (u *listByCourseUsecase) Execute(ctx context.Context, courseID string) (*CourseLearningStrategyListOutput, apperrors.ApplicationError) {
	app := u.factory()

	strategies, err := app.Repositories.LearningStrategy.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyListError, err)
	}

	data := make([]CourseLearningStrategyData, 0, len(strategies))
	for _, cls := range strategies {
		data = append(data, toCourseStrategyData(cls))
	}

	return &CourseLearningStrategyListOutput{Data: data}, nil
}
