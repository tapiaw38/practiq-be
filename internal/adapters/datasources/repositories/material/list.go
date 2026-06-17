package material

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) List(ctx context.Context, filter ListFilter) ([]domain.Material, error) {
	query := `
		SELECT m.id, m.course_id, m.teacher_id, m.title, m.type, COALESCE(m.file_url,''), COALESCE(m.extracted_text,''), m.status, m.created_at
		FROM materials m
		JOIN courses c ON c.id = m.course_id
		WHERE m.course_id = $1 AND c.deleted_at IS NULL
		ORDER BY m.created_at DESC`

	args := []interface{}{filter.CourseID}
	argIndex := 2

	if filter.Limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(` OFFSET $%d`, argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
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
