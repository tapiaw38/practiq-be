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
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Title            string   `json:"title" binding:"required"`
		TopicID          string   `json:"topic_id"`
		Level            int      `json:"level"`
		SheetType        string   `json:"sheet_type"`
		TestStyle        string   `json:"test_style"`
		ExerciseIDs      []string `json:"exercise_ids"`
		ScheduledAt      string   `json:"scheduled_at"`
		AvailableUntil   string   `json:"available_until"`
		MaxAttempts      *int     `json:"max_attempts"`
		TimeLimitMinutes *int     `json:"time_limit_minutes"`
	}

	UpdateOutput struct {
		Data PracticeSheetData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	if appErr := validateSheetTypeAndTestStyle(input.SheetType, input.TestStyle); appErr != nil {
		return nil, appErr
	}
	scheduledAt, availableUntil, appErr := resolveWindow(input.ScheduledAt, input.AvailableUntil)
	if appErr != nil {
		return nil, appErr
	}
	maxAttempts := positiveOrNil(input.MaxAttempts)
	timeLimitMinutes := positiveOrNil(input.TimeLimitMinutes)
	app := u.contextFactory()

	ps, err := app.Repositories.PracticeSheet.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, ps.CourseID); appErr != nil {
		return nil, appErr
	}

	level := input.Level
	if level < 1 {
		level = 1
	}
	if input.SheetType == sheetTypeLevelTest {
		exists, err := app.Repositories.PracticeSheet.HasOtherLevelTest(ctx, ps.CourseID, level, id)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetUpdateError, err)
		}
		if exists {
			return nil, apperrors.NewBadRequestError("this level already has a level test; edit or delete it first")
		}
	}
	if input.SheetType != sheetTypeLevelTest {
		scheduledAt = nil
	}

	if err := app.Repositories.PracticeSheet.Update(ctx, id, domain.PracticeSheet{
		Title:            input.Title,
		TopicID:          input.TopicID,
		Level:            level,
		SheetType:        input.SheetType,
		TestStyle:        input.TestStyle,
		ScheduledAt:      scheduledAt,
		MaxAttempts:      maxAttempts,
		TimeLimitMinutes: timeLimitMinutes,
		AvailableUntil:   availableUntil,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetUpdateError, err)
	}

	if err := app.Repositories.PracticeSheet.ReplaceExercises(ctx, id, input.ExerciseIDs); err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetUpdateError, err)
	}

	ps, err = app.Repositories.PracticeSheet.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	notifyScheduledLevelTest(ctx, app, *ps)

	return &UpdateOutput{Data: toSheetData(app, *ps, true)}, nil
}
