package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	CreateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, in CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		TopicID          string   `json:"topic_id"`
		StrategyID       string   `json:"strategy_id"`
		Title            string   `json:"title" binding:"required"`
		Level            int      `json:"level"`
		SheetType        string   `json:"sheet_type"`
		TestStyle        string   `json:"test_style"`
		ExerciseIDs      []string `json:"exercise_ids"`
		ScheduledAt      string   `json:"scheduled_at"`
		AvailableUntil   string   `json:"available_until"`
		MaxAttempts      *int     `json:"max_attempts"`
		TimeLimitMinutes *int     `json:"time_limit_minutes"`
	}

	CreateOutput struct {
		Data PracticeSheetData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, in CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	if appErr := validateSheetTypeAndTestStyle(in.SheetType, in.TestStyle); appErr != nil {
		return nil, appErr
	}
	scheduledAt, availableUntil, appErr := resolveWindow(in.ScheduledAt, in.AvailableUntil)
	if appErr != nil {
		return nil, appErr
	}
	maxAttempts := positiveOrNil(in.MaxAttempts)
	timeLimitMinutes := positiveOrNil(in.TimeLimitMinutes)
	app := u.contextFactory()

	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
	}

	level := in.Level
	if level < 1 {
		level = 1
	}

	sheetType := in.SheetType
	if sheetType != "level_test" {
		sheetType = "practice"
	}
	if sheetType == sheetTypeLevelTest {
		exists, err := app.Repositories.PracticeSheet.HasOtherLevelTest(ctx, courseID, level, "")
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetCreateError, err)
		}
		if exists {
			return nil, apperrors.NewBadRequestError("this level already has a level test; edit or delete it first")
		}
	}
	testStyle := in.TestStyle
	if testStyle != "canvas" {
		testStyle = "keyboard"
	}

	if sheetType != sheetTypeLevelTest {
		scheduledAt = nil
	}
	id, err := app.Repositories.PracticeSheet.Create(ctx, domain.PracticeSheet{
		CourseID:         courseID,
		TopicID:          in.TopicID,
		StrategyID:       in.StrategyID,
		Title:            in.Title,
		Level:            level,
		SheetType:        sheetType,
		TestStyle:        testStyle,
		ScheduledAt:      scheduledAt,
		MaxAttempts:      maxAttempts,
		TimeLimitMinutes: timeLimitMinutes,
		AvailableUntil:   availableUntil,
		CreatedBy:        "teacher",
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetCreateError, err)
	}

	for i, exerciseID := range in.ExerciseIDs {
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
