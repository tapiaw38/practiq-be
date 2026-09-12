package course

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Update(ctx context.Context, id string, c domain.Course) error {
	query := `UPDATE courses SET title=$1, description=$2, level=$3, subject=$4, grade_id=$5, subject_id=$6 WHERE id=$7`
	_, err := r.db.ExecContext(ctx, query, c.Title, c.Description, c.Level, c.Subject, nullableUUID(c.GradeID), nullableUUID(c.SubjectID), id)
	return err
}
