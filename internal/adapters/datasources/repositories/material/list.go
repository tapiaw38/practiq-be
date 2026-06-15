package material

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) List(ctx context.Context, courseID string) ([]domain.Material, error) {
	query := `SELECT id, course_id, teacher_id, title, type, COALESCE(file_url,''), COALESCE(extracted_text,''), status, created_at FROM materials WHERE course_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []domain.Material
	for rows.Next() {
		var m domain.Material
		if err := rows.Scan(&m.ID, &m.CourseID, &m.TeacherID, &m.Title, &m.Type, &m.FileURL, &m.ExtractedText, &m.Status, &m.CreatedAt); err != nil {
			return nil, err
		}
		materials = append(materials, m)
	}
	return materials, nil
}
