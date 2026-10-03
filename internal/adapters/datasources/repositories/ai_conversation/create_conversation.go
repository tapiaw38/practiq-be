package aiconversation

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) CreateConversation(ctx context.Context, c domain.AIConversation) (string, error) {
	query := `
		INSERT INTO ai_conversations (student_id, course_id, practice_sheet_id)
		VALUES ($1, NULLIF($2,'')::uuid, NULLIF($3,'')::uuid)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, c.StudentID, c.CourseID, c.PracticeSheetID).Scan(&id)
	return id, err
}
