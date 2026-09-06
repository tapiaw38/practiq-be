package course

import (
	"context"

	"github.com/lib/pq"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

// ListDashboardSummaries returns, for every course a student is enrolled in,
// the counts and level the home screen shows.
//
// Set-based on purpose. The screen used to ask for practice sheets, notebooks
// and levels once per course, which was three round trips per course on top of
// the courses call itself. Rewriting that as a loop of queries here would move
// the N+1 from the network to the database rather than remove it, so each fact
// is one aggregate over all the student's courses at once.
func (r *repository) ListDashboardSummaries(ctx context.Context, studentID string) ([]domain.CourseDashboardSummary, error) {
	rows, err := r.db.QueryContext(ctx, `
		WITH student_courses AS (
			SELECT c.id, c.title, COALESCE(c.subject, '') AS subject, c.created_at
			FROM courses c
			WHERE c.deleted_at IS NULL
			  AND ($2 = '' OR c.school_id = NULLIF($2, '')::uuid)
			  -- A student reaches a course either by enrolling in it directly or
			  -- by belonging to its grade, and the grade is the usual route.
			  -- Matching on enrolments alone returned an empty home for those
			  -- students. Same rule the course listing applies.
			  AND (
			    EXISTS (SELECT 1 FROM enrollments e
			            WHERE e.course_id = c.id AND e.student_id = $1)
			    OR EXISTS (SELECT 1 FROM grade_memberships gm
			               WHERE gm.grade_id = c.grade_id AND gm.user_id = $1)
			  )
		),
		sheet_counts AS (
			SELECT ps.course_id,
			       COUNT(*) FILTER (WHERE COALESCE(ps.sheet_type,'practice') <> 'level_test') AS practices,
			       COUNT(*) FILTER (WHERE ps.sheet_type = 'level_test') AS level_tests
			FROM practice_sheets ps
			JOIN student_courses sc ON sc.id = ps.course_id
			GROUP BY ps.course_id
		),
		notebook_counts AS (
			SELECT n.course_id, COUNT(*) AS notebooks
			FROM notebooks n
			JOIN student_courses sc ON sc.id = n.course_id
			WHERE n.deleted_at IS NULL
			GROUP BY n.course_id
		),
		course_topics AS (
			SELECT t.course_id, ARRAY_AGG(DISTINCT t.id::text) AS topic_ids
			FROM topics t
			JOIN student_courses sc ON sc.id = t.course_id
			GROUP BY t.course_id
		),
		levels AS (
			SELECT cp.course_id, cp.current_level
			FROM student_course_progress cp
			WHERE cp.student_id = $1
		)
		SELECT sc.id, sc.title, sc.subject,
		       COALESCE(s.practices, 0), COALESCE(s.level_tests, 0),
		       COALESCE(nb.notebooks, 0), COALESCE(l.current_level, 1),
		       COALESCE(ct.topic_ids, ARRAY[]::text[])
		FROM student_courses sc
		LEFT JOIN sheet_counts s ON s.course_id = sc.id
		LEFT JOIN notebook_counts nb ON nb.course_id = sc.id
		LEFT JOIN levels l ON l.course_id = sc.id
		LEFT JOIN course_topics ct ON ct.course_id = sc.id
		ORDER BY sc.created_at DESC
	`, studentID, tenant.SchoolID(ctx))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries := []domain.CourseDashboardSummary{}
	for rows.Next() {
		var s domain.CourseDashboardSummary
		if err := rows.Scan(&s.CourseID, &s.Title, &s.Subject,
			&s.PracticeSheets, &s.LevelTests, &s.Notebooks, &s.CurrentLevel, pq.Array(&s.TopicIDs)); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}
