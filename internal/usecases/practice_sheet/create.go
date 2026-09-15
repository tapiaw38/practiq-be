package practicesheet

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	CreateUsecase interface {
		Execute(context.Context, string, bool, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		CourseID       string
		TopicID        string `json:"topic_id"`
		StrategyID     string `json:"strategy_id"`
		Title          string `json:"title"`
		Level          int    `json:"level"`
		SheetType      string `json:"sheet_type"`
		TestStyle      string `json:"test_style"`
		ScheduledAt    *time.Time
		AvailableUntil *time.Time
		ExerciseIDs    []string `json:"exercise_ids"`
	}

	CreateOutput struct {
		Data PracticeSheetData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, input.CourseID); appErr != nil {
		return nil, appErr
	}

	level := input.Level
	if level < 1 {
		level = 1
	}

	sheetType := input.SheetType
	if sheetType != "level_test" {
		sheetType = "practice"
	}
	if sheetType == sheetTypeLevelTest {
		exists, err := app.Repositories.PracticeSheet.HasOtherLevelTest(ctx, input.CourseID, level, "")
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetCreateError, err)
		}
		if exists {
			return nil, apperrors.NewBadRequestError("this level already has a level test; edit or delete it first")
		}
	}
	testStyle := input.TestStyle
	if testStyle != "canvas" {
		testStyle = "keyboard"
	}
	// Only level tests are scheduled; a date on a practice sheet would lock it
	// for students with no way to see why.
	scheduledAt := input.ScheduledAt
	if sheetType != sheetTypeLevelTest {
		scheduledAt = nil
	}
	id, err := app.Repositories.PracticeSheet.Create(ctx, domain.PracticeSheet{
		CourseID:       input.CourseID,
		TopicID:        input.TopicID,
		StrategyID:     input.StrategyID,
		Title:          input.Title,
		Level:          level,
		SheetType:      sheetType,
		TestStyle:      testStyle,
		ScheduledAt:    scheduledAt,
		AvailableUntil: input.AvailableUntil,
		CreatedBy:      "teacher",
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

	notifyScheduledLevelTest(ctx, app, *ps)

	return &CreateOutput{Data: toSheetData(app, *ps, true)}, nil
}
