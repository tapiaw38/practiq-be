package studentcoursexp

import "context"

func (r *repository) ListByStudent(ctx context.Context, studentID string) ([]CourseXP, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT course_id::text, total_xp
		FROM student_course_xp
		WHERE student_id = $1`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]CourseXP, 0)
	for rows.Next() {
		var item CourseXP
		if err := rows.Scan(&item.CourseID, &item.TotalXP); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
