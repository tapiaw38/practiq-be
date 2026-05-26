package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type UpdateUsecase interface {
	Execute(context.Context, string, UpdateInput) (*PracticeSheetOutput, apperrors.ApplicationError)
}

type updateUsecase struct {
	factory appcontext.Factory
}

type UpdateInput struct {
	Title       string
	TopicID     string
	Level       int
	SheetType   string
	TestStyle   string
	ExerciseIDs []string
}

func NewUpdateUsecase(factory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{factory: factory}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*PracticeSheetOutput, apperrors.ApplicationError) {
	app := u.factory()

	if err := app.Repositories.PracticeSheet.Update(ctx, id, domain.PracticeSheet{
		Title:     input.Title,
		TopicID:   input.TopicID,
		Level:     input.Level,
		SheetType: input.SheetType,
		TestStyle: input.TestStyle,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetUpdateError, err)
	}

	if err := app.Repositories.PracticeSheet.ReplaceExercises(ctx, id, input.ExerciseIDs); err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetUpdateError, err)
	}

	ps, err := app.Repositories.PracticeSheet.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	return &PracticeSheetOutput{Data: toSheetData(*ps)}, nil
}
