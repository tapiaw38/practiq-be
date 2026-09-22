package studentcoursexp

import "context"

// LeaderboardByCourse returns the top of the course plus the given student's
// own row, even when it falls outside the top. Ranking happens in the database
// so the position is right whatever slice we ship: counting in Go over a
// truncated list would put everyone outside the top at the same place.
//
// Ties share a position (1, 1, 3) and break by name only to keep the order
// stable between calls. A student with no XP yet has no row here at all; the
// caller decides what to show them.
func (r *repository) LeaderboardByCourse(ctx context.Context, courseID, studentID string, top int) ([]LeaderboardEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		WITH ranked AS (
			SELECT
				x.student_id,
				COALESCE(p.name, '') AS name,
				x.total_xp,
				RANK() OVER (ORDER BY x.total_xp DESC, COALESCE(p.name, '') ASC) AS position
			FROM student_course_xp x
			JOIN user_profiles p ON p.id = x.student_id
			WHERE x.course_id = $1::uuid
		)
		SELECT student_id, name, total_xp, position
		FROM ranked
		WHERE position <= $2 OR student_id = $3
		ORDER BY position`, courseID, top, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]LeaderboardEntry, 0, top+1)
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.StudentID, &entry.Name, &entry.TotalXP, &entry.Position); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}
