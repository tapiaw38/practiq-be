package studentattempt

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetDailyAttempts(ctx context.Context, studentID, courseID string, from, to *time.Time) ([]domain.DailyAttemptCount, error) {
	query := `
		SELECT
			DATE(sa.created_at) as date,
			COUNT(*) as total,
			SUM(CASE WHEN sa.score >= 100 THEN 1 ELSE 0 END) as correct
		FROM student_attempts sa
	`
	args := []interface{}{studentID}
	argIndex := 2

	if courseID != "" {
		query += `
		JOIN exercises e ON e.id = sa.exercise_id
		JOIN topics t ON t.id = e.topic_id
		WHERE sa.student_id = $1 AND t.course_id = $` + itoa(argIndex)
		args = append(args, courseID)
		argIndex++
	} else {
		query += ` WHERE sa.student_id = $1`
	}

	if from != nil {
		query += ` AND sa.created_at >= $` + itoa(argIndex)
		args = append(args, *from)
		argIndex++
	}

	if to != nil {
		query += ` AND sa.created_at <= $` + itoa(argIndex)
		args = append(args, *to)
		argIndex++
	}

	query += `
		GROUP BY DATE(sa.created_at)
		ORDER BY DATE(sa.created_at) DESC
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.DailyAttemptCount
	for rows.Next() {
		var d domain.DailyAttemptCount
		if err := rows.Scan(&d.Date, &d.Total, &d.Correct); err != nil {
			return nil, err
		}
		result = append(result, d)
	}

	return result, nil
}

func itoa(n int) string {
	if n < 10 {
		return string('0' + byte(n))
	}
	return string('0'+byte(n/10)) + string('0'+byte(n%10))
}
