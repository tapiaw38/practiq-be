package learningstrategy

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListByCourseUsecase interface {
		Execute(context.Context, string) (*ListByCourseOutput, apperrors.ApplicationError)
	}

	listByCourseUsecase struct {
		contextFactory appcontext.Factory
	}

	ListByCourseOutput struct {
		Data []CourseLearningStrategyData `json:"data"`
	}
)

func NewListByCourseUsecase(contextFactory appcontext.Factory) ListByCourseUsecase {
	return &listByCourseUsecase{contextFactory: contextFactory}
}

func (u *listByCourseUsecase) Execute(ctx context.Context, courseID string) (*ListByCourseOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	strategies, err := app.Repositories.LearningStrategy.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyListError, err)
	}

	data := make([]CourseLearningStrategyData, 0, len(strategies))
	for _, cls := range strategies {
		data = append(data, toCourseStrategyData(cls))
	}

	return &ListByCourseOutput{Data: data}, nil
}
