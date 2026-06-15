package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
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
		SET title = $1, topic_id = NULLIF($2,'')::uuid, level = $3, sheet_type = $4, test_style = $5
		WHERE id = $6
	`, ps.Title, ps.TopicID, level, sheetType, testStyle, id)
	return err
}
