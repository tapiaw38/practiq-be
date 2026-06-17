package enrollment

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListStudents(ctx context.Context, filter ListFilter) ([]domain.UserProfile, error) {
	query := `
		SELECT DISTINCT up.id, up.name, up.email, up.profile_type, up.created_at
		FROM user_profiles up
		JOIN courses c ON c.id = $1
		LEFT JOIN enrollments e ON e.course_id = c.id AND e.student_id = up.id
		LEFT JOIN grade_memberships gm ON gm.grade_id = c.grade_id AND gm.user_id = up.id
		WHERE up.profile_type = 'student'
		  AND (e.student_id IS NOT NULL OR gm.user_id IS NOT NULL)
		ORDER BY up.name ASC`

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

	var students []domain.UserProfile
	for rows.Next() {
		var s domain.UserProfile
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.ProfileType, &s.CreatedAt); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}
