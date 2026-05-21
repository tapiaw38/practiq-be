package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type CreateUsecase interface {
	Execute(context.Context, CreateInput) (*PracticeSheetOutput, apperrors.ApplicationError)
}

type createUsecase struct {
	factory appcontext.Factory
}

type CreateInput struct {
	CourseID    string
	TopicID     string   `json:"topic_id"`
	StrategyID  string   `json:"strategy_id"`
	Title       string   `json:"title"`
	Level       int      `json:"level"`
	SheetType   string   `json:"sheet_type"`
	TestStyle   string   `json:"test_style"`
	ExerciseIDs []string `json:"exercise_ids"`
}

func NewCreateUsecase(factory appcontext.Factory) CreateUsecase {
	return &createUsecase{factory: factory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*PracticeSheetOutput, apperrors.ApplicationError) {
	app := u.factory()

	level := input.Level
	if level < 1 {
		level = 1
	}

	sheetType := input.SheetType
	if sheetType != "level_test" {
		sheetType = "practice"
	}
	testStyle := input.TestStyle
	if testStyle != "canvas" {
		testStyle = "keyboard"
	}
	id, err := app.Repositories.PracticeSheet.Create(ctx, domain.PracticeSheet{
		CourseID:   input.CourseID,
		TopicID:    input.TopicID,
		StrategyID: input.StrategyID,
		Title:      input.Title,
		Level:      level,
		SheetType:  sheetType,
		TestStyle:  testStyle,
		CreatedBy:  "teacher",
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetCreateError, err)
	}

	for i, exerciseID := range input.ExerciseIDs {
		if err := app.Repositories.PracticeSheet.AddExercise(ctx, id, exerciseID, i); err != nil {
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetCreateError, err)
		}
	}

	ps, err := app.Repositories.PracticeSheet.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}

	return &PracticeSheetOutput{Data: toSheetData(*ps)}, nil
}
