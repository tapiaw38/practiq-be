package material

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Material, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT m.id, m.course_id, m.teacher_id, m.title, m.type, COALESCE(m.file_url,''), COALESCE(m.extracted_text,''), m.status, m.created_at
		FROM materials m
		JOIN courses c ON c.id = m.course_id
		WHERE m.id = $1 AND c.deleted_at IS NULL
	`, id)

	var m domain.Material
	if err := row.Scan(&m.ID, &m.CourseID, &m.TeacherID, &m.Title, &m.Type, &m.FileURL, &m.ExtractedText, &m.Status, &m.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
