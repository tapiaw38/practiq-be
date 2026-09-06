package material

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

// List returns the materials of a course with ExtractedText cut to a preview.
//
// The column holds the whole extracted document. The listing shows two clamped
// lines of it, so pulling the rest moves the entire text of every material in
// the course over the wire to be thrown away. Whoever needs all of it reads one
// material through Get.
//
// The cut is 501 characters so the caller can tell a text that ends exactly at
// the preview length from one that continues.
func (r *repository) List(ctx context.Context, filter ListFilter) ([]domain.Material, error) {
	query := `
		SELECT m.id, m.course_id, m.teacher_id, m.title, m.type, COALESCE(m.file_url,''), left(COALESCE(m.extracted_text,''), 501), m.status, m.created_at
		FROM materials m
		JOIN courses c ON c.id = m.course_id
		WHERE m.course_id = $1 AND c.deleted_at IS NULL AND ($2 = '' OR c.school_id = NULLIF($2, '')::uuid)
		ORDER BY m.created_at DESC`

	args := []interface{}{filter.CourseID, tenant.SchoolID(ctx)}
	argIndex := 3

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
