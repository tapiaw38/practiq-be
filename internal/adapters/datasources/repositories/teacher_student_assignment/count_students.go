package teacherstudentassignment

import "context"

// CountStudents counts the distinct students a teacher is responsible for.
//
// A student reaches a teacher by two routes, and both have to count or the
// number on the subscription tab will disagree with the student list the
// teacher is looking at. The same split already produced a bug once, when a
// query matched enrolments alone and a student's home came back empty.
//
// Counted here rather than derived from any single list because the two routes
// overlap: a student assigned to a teacher who is also in one of that teacher's
// grades is one student, not two.
func (r *repository) CountStudents(ctx context.Context, teacherID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM (
			SELECT tsa.student_id AS id
			FROM teacher_student_assignments tsa
			WHERE tsa.teacher_id = $1 AND tsa.status = 'active'

			UNION

			SELECT e.student_id
			FROM enrollments e
			JOIN courses c ON c.id = e.course_id
			WHERE c.teacher_id = $1 AND c.deleted_at IS NULL

			UNION

			SELECT gm.user_id
			FROM grade_memberships gm
			JOIN courses c ON c.grade_id = gm.grade_id
			JOIN user_profiles up ON up.id = gm.user_id
			WHERE c.teacher_id = $1 AND c.deleted_at IS NULL
			  AND up.profile_type = 'student'
		) students
	`, teacherID).Scan(&count)
	return count, err
}
