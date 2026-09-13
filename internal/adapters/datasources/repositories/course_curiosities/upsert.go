package coursecuriosities

import (
	"context"
	"encoding/json"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Upsert(ctx context.Context, c domain.CourseCuriosities) error {
	curiositiesJSON := "[]"
	if len(c.Curiosities) > 0 {
		if b, err := json.Marshal(c.Curiosities); err == nil {
			curiositiesJSON = string(b)
		}
	}

	query := `
		INSERT INTO course_curiosities (course_id, curiosities)
		VALUES ($1, $2)
		ON CONFLICT (course_id) DO UPDATE SET curiosities = $2, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, c.CourseID, curiositiesJSON)
	return err
}
