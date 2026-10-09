package school

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

const activeStudentsQuery = `
	FROM school_members sm
	LEFT JOIN (
		SELECT student_id, MAX(last_practiced_at) AS last_at
		FROM student_topic_progress GROUP BY student_id
	) activity ON activity.student_id = sm.user_id
	WHERE sm.school_id = $1 AND sm.role = 'student' AND sm.active
	ORDER BY activity.last_at ASC NULLS FIRST, sm.created_at ASC
`

func (r *repository) ListStudentsWithActivity(ctx context.Context, schoolID string) ([]domain.StudentActivity, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT sm.user_id, activity.last_at "+activeStudentsQuery, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.StudentActivity{}
	for rows.Next() {
		var (
			id     string
			lastAt sql.NullTime
		)
		if err := rows.Scan(&id, &lastAt); err != nil {
			return nil, err
		}
		student := domain.StudentActivity{UserID: id}
		if lastAt.Valid {
			at := lastAt.Time
			student.LastPracticed = &at
		}
		out = append(out, student)
	}
	return out, rows.Err()
}

func (r *repository) ListStudentsByActivity(ctx context.Context, schoolID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT sm.user_id "+activeStudentsQuery, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
