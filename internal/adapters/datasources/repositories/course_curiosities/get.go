package coursecuriosities

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, courseID string) (*domain.CourseCuriosities, error) {
	query := `SELECT course_id, curiosities FROM course_curiosities WHERE course_id = $1`
	row := r.db.QueryRowContext(ctx, query, courseID)

	var c domain.CourseCuriosities
	var curiositiesJSON string
	err := row.Scan(&c.CourseID, &curiositiesJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if curiositiesJSON != "" {
		_ = json.Unmarshal([]byte(curiositiesJSON), &c.Curiosities)
	}
	return &c, nil
}
