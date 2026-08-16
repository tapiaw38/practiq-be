package studentattempt

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// ListPendingReview returns answers requiring a teacher from that teacher's
// courses. Statement image/audio answers are included because text-only
// automatic grading cannot assess them safely.
func (r *repository) ListPendingReview(ctx context.Context, teacherID string, includeReviewed bool) ([]domain.PendingAttemptReview, error) {
	query := `
		SELECT sa.id, sa.student_id, COALESCE(up.name, ''), sa.exercise_id, e.question, e.type, COALESCE(e.metadata::text, ''),
		       COALESCE(sa.practice_sheet_id::text, ''), COALESCE(ps.title, ''), COALESCE(ps.sheet_type, ''),
		       c.id, c.title,
		       COALESCE(sa.image_url, ''), COALESCE(sa.attachment_url, ''), COALESCE(sa.attachment_name, ''), COALESCE(sa.attachment_content_type, ''),
		       COALESCE(sa.answer_text, ''), COALESCE(sa.ai_feedback, ''), sa.ai_is_correct,
		       sa.teacher_is_correct, COALESCE(sa.teacher_feedback, ''), sa.teacher_reviewed_at, sa.created_at
		FROM student_attempts sa
		JOIN exercises e ON e.id = sa.exercise_id
		LEFT JOIN practice_sheets ps ON ps.id = sa.practice_sheet_id
		JOIN topics t ON t.id = e.topic_id
		JOIN courses c ON c.id = t.course_id
		LEFT JOIN user_profiles up ON up.id = sa.student_id
		WHERE (COALESCE(sa.attachment_url, '') <> '' OR sa.needs_teacher_review)
		  AND c.teacher_id = $1
		  AND c.deleted_at IS NULL
	`
	if !includeReviewed {
		query += ` AND sa.needs_teacher_review AND sa.teacher_reviewed_at IS NULL`
	}
	query += ` ORDER BY sa.created_at DESC LIMIT 100`

	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []domain.PendingAttemptReview{}
	for rows.Next() {
		var review domain.PendingAttemptReview
		var exerciseMetadata string
		if err := rows.Scan(
			&review.AttemptID,
			&review.StudentID,
			&review.StudentName,
			&review.ExerciseID,
			&review.Question,
			&review.ExerciseType,
			&exerciseMetadata,
			&review.PracticeSheetID,
			&review.PracticeSheetTitle,
			&review.SheetType,
			&review.CourseID,
			&review.CourseTitle,
			&review.ImageURL,
			&review.AttachmentURL,
			&review.AttachmentName,
			&review.AttachmentContentType,
			&review.AnswerText,
			&review.AIFeedback,
			&review.AIIsCorrect,
			&review.TeacherIsCorrect,
			&review.TeacherFeedback,
			&review.TeacherReviewedAt,
			&review.CreatedAt,
		); err != nil {
			return nil, err
		}
		exercise := domain.Exercise{Metadata: exerciseMetadata}
		review.StatementMediaURL = exercise.MediaURL()
		// Only the flag: the image itself is fetched when the teacher opens it.
		review.HasTeacherImage = exercise.TeacherImage() != ""
		reviews = append(reviews, review)
	}
	return reviews, rows.Err()
}

// GetTeacherForAttempt resolves who owns the course an attempt belongs to, so
// the review endpoint can reject anyone else.
func (r *repository) GetTeacherForAttempt(ctx context.Context, attemptID string) (string, error) {
	var teacherID string
	err := r.db.QueryRowContext(ctx, `
		SELECT c.teacher_id
		FROM student_attempts sa
		JOIN exercises e ON e.id = sa.exercise_id
		JOIN topics t ON t.id = e.topic_id
		JOIN courses c ON c.id = t.course_id
		WHERE sa.id = $1
	`, attemptID).Scan(&teacherID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return teacherID, err
}

// GetExerciseIDForAttempt resolves which exercise an attempt answered, so the
// statement can be fetched without trusting an id from the client.
func (r *repository) GetExerciseIDForAttempt(ctx context.Context, attemptID string) (string, error) {
	var exerciseID string
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(exercise_id::text, '') FROM student_attempts WHERE id = $1
	`, attemptID).Scan(&exerciseID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return exerciseID, err
}

// Review records the teacher's verdict. The score follows the verdict so
// progress reflects the corrected answer.
func (r *repository) Review(ctx context.Context, attemptID string, isCorrect bool, feedback string) error {
	score := 0.0
	if isCorrect {
		score = 100.0
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE student_attempts
		SET teacher_is_correct = $1,
		    teacher_feedback = $2,
		    teacher_reviewed_at = NOW(),
		    is_correct = $1,
		    score = $3
		WHERE id = $4
	`, isCorrect, feedback, score, attemptID)
	return err
}
