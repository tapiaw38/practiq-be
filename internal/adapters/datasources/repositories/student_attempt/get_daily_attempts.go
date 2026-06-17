package studentattempt

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetDailyAttempts(ctx context.Context, studentID string, from, to *time.Time) ([]domain.DailyAttemptCount, error) {
	query := `
		SELECT
			DATE(created_at) as date,
			COUNT(*) as total,
			SUM(CASE WHEN score >= 100 THEN 1 ELSE 0 END) as correct
		FROM student_attempts
		WHERE student_id = $1
	`
	args := []interface{}{studentID}
	argIndex := 2

	if from != nil {
		query += ` AND created_at >= $` + itoa(argIndex)
		args = append(args, *from)
		argIndex++
	}

	if to != nil {
		query += ` AND created_at <= $` + itoa(argIndex)
		args = append(args, *to)
		argIndex++
	}

	query += `
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at) DESC
		LIMIT 30
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
