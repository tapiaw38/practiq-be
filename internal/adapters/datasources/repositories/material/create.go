package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Create(ctx context.Context, m domain.Material) (string, error) {
	query := `
		INSERT INTO materials (course_id, teacher_id, title, type, file_url, extracted_text, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, m.CourseID, m.TeacherID, m.Title, m.Type, m.FileURL, m.ExtractedText, m.Status).Scan(&id)
	return id, err
}
