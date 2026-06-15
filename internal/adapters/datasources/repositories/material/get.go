package material

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.Material, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, course_id, teacher_id, title, type, COALESCE(file_url,''), COALESCE(extracted_text,''), status, created_at
		FROM materials WHERE id = $1
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
