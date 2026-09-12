package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, ps domain.PracticeSheet) (string, error) {
	sheetType := ps.SheetType
	if sheetType != "level_test" {
		sheetType = "practice"
	}
	testStyle := ps.TestStyle
	if testStyle != "canvas" {
		testStyle = "keyboard"
	}
	query := `
		INSERT INTO practice_sheets (course_id, topic_id, strategy_id, title, level, sheet_type, test_style, created_by, scheduled_at, available_until)
		VALUES ($1, NULLIF($2,'')::uuid, NULLIF($3,'')::uuid, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, ps.CourseID, ps.TopicID, ps.StrategyID, ps.Title, ps.Level, sheetType, testStyle, ps.CreatedBy, ps.ScheduledAt, ps.AvailableUntil).Scan(&id)
	return id, err
}
