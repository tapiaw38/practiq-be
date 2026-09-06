package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) Update(ctx context.Context, id string, ps domain.PracticeSheet) error {
	sheetType := ps.SheetType
	if sheetType != "level_test" {
		sheetType = "practice"
	}
	testStyle := ps.TestStyle
	if testStyle != "canvas" {
		testStyle = "keyboard"
	}
	level := ps.Level
	if level < 1 {
		level = 1
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE practice_sheets
		SET title = $1, topic_id = NULLIF($2,'')::uuid, level = $3, sheet_type = $4, test_style = $5, scheduled_at = $6, available_until = $7
		WHERE id = $8 AND ($9 = '' OR school_id = NULLIF($9, '')::uuid)
	`, ps.Title, ps.TopicID, level, sheetType, testStyle, ps.ScheduledAt, ps.AvailableUntil, id, tenant.SchoolID(ctx))
	return err
}
