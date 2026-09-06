package aiconversation

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

// ListRecentHelpRequests keeps exercise memory bounded and scoped to its owner.
func (r *repository) ListRecentHelpRequests(ctx context.Context, studentID, exerciseID string, limit int) ([]domain.AIHelpRequest, error) {
	query := `
		SELECT id, student_id, COALESCE(exercise_id::text, ''), question,
			COALESCE(ai_response, ''), COALESCE(help_type, ''), created_at
		FROM ai_help_requests
		WHERE student_id = $1 AND exercise_id = NULLIF($2, '')::uuid AND ($4 = '' OR school_id = NULLIF($4, '')::uuid)
		ORDER BY created_at DESC
		LIMIT $3
	`
	rows, err := r.db.QueryContext(ctx, query, studentID, exerciseID, limit, tenant.SchoolID(ctx))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.AIHelpRequest, 0)
	for rows.Next() {
		var item domain.AIHelpRequest
		if err := rows.Scan(&item.ID, &item.StudentID, &item.ExerciseID, &item.Question, &item.AIResponse, &item.HelpType, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
