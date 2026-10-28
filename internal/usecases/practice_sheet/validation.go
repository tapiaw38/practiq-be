package practicesheet

import (
	"strings"
	"time"

	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

func validateSheetTypeAndTestStyle(sheetType, testStyle string) apperrors.ApplicationError {
	if sheetType != "" && sheetType != "practice" && sheetType != "level_test" {
		return apperrors.NewBadRequestError("sheet_type must be 'practice' or 'level_test'")
	}
	if testStyle != "" && testStyle != "keyboard" && testStyle != "canvas" {
		return apperrors.NewBadRequestError("test_style must be 'keyboard' or 'canvas'")
	}
	return nil
}

func parseScheduledAt(value string) (*time.Time, apperrors.ApplicationError) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, apperrors.NewBadRequestError("scheduled_at must be an RFC 3339 timestamp, e.g. 2026-09-01T14:30:00Z")
	}
	utc := parsed.UTC()
	return &utc, nil
}

func invalidWindow(message string) apperrors.ApplicationError {
	return apperrors.NewApplicationError(mappings.ErrorDetails{
		InternalCode: "practice_sheet:invalid-window",
		StatusCode:   400,
		Message:      message,
	}, nil)
}

func resolveWindow(scheduledAtRaw, availableUntilRaw string) (*time.Time, *time.Time, apperrors.ApplicationError) {
	scheduledAt, appErr := parseScheduledAt(scheduledAtRaw)
	if appErr != nil {
		return nil, nil, appErr
	}
	availableUntil, appErr := parseScheduledAt(availableUntilRaw)
	if appErr != nil {
		return nil, nil, appErr
	}
	if availableUntil != nil && scheduledAt == nil {
		return nil, nil, invalidWindow("the closing date requires an opening date")
	}
	if scheduledAt != nil && availableUntil != nil && !availableUntil.After(*scheduledAt) {
		return nil, nil, invalidWindow("the closing date must be after the opening date")
	}
	return scheduledAt, availableUntil, nil
}

func positiveOrNil(value *int) *int {
	if value == nil || *value <= 0 {
		return nil
	}
	return value
}
