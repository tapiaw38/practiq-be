package aiconversation

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.AIConversation, error) {
	query := `
		SELECT id, student_id, COALESCE(course_id::text, ''), COALESCE(practice_sheet_id::text, ''), created_at
		FROM ai_conversations
		WHERE id = $1
	`
	var c domain.AIConversation
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.StudentID,
		&c.CourseID,
		&c.PracticeSheetID,
		&c.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
