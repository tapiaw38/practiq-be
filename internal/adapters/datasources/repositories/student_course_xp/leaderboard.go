package studentcoursexp

import "context"

// LeaderboardByCourse ranks everyone who belongs to the course, whether they
// have earned anything or not: a student with no XP row is on zero, not
// missing. Membership is either a direct enrolment or the course's grade,
// the same pair every other course read resolves against.
//
// Two window functions rather than one. Position is a RANK so equal scores
// share a place — at the start of a course everybody is on zero and printing
// 1, 2, 3 down an arbitrary order would invent a ranking nobody earned. The
// cut to the top N rides on ROW_NUMBER instead, because a RANK of 1 shared by
// thirty students would let all thirty through a `position <= 10` filter.
func (r *repository) LeaderboardByCourse(ctx context.Context, courseID, studentID string, top int) ([]LeaderboardEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		WITH members AS (
			SELECT e.student_id AS student_id
			FROM enrollments e
			WHERE e.course_id = $1::uuid
			UNION
			SELECT gm.user_id AS student_id
			FROM grade_memberships gm
			JOIN courses c ON c.grade_id = gm.grade_id
			WHERE c.id = $1::uuid
		),
		ranked AS (
			SELECT
				m.student_id,
				COALESCE(x.total_xp, 0) AS total_xp,
				RANK() OVER (ORDER BY COALESCE(x.total_xp, 0) DESC) AS position,
				ROW_NUMBER() OVER (ORDER BY COALESCE(x.total_xp, 0) DESC, m.student_id) AS seq
			FROM members m
			LEFT JOIN student_course_xp x
				ON x.student_id = m.student_id AND x.course_id = $1::uuid
		)
		SELECT student_id, total_xp, position
		FROM ranked
		WHERE seq <= $2 OR student_id = $3
		ORDER BY seq`, courseID, top, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]LeaderboardEntry, 0, top+1)
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.StudentID, &entry.TotalXP, &entry.Position); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}
