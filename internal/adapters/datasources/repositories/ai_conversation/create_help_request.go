package aiconversation

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) CreateHelpRequest(ctx context.Context, h domain.AIHelpRequest) (string, error) {
	query := `
		INSERT INTO ai_help_requests (student_id, exercise_id, question, ai_response, help_type)
		VALUES ($1, NULLIF($2,'')::uuid, $3, $4, $5)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, h.StudentID, h.ExerciseID, h.Question, h.AIResponse, h.HelpType).Scan(&id)
	return id, err
}
