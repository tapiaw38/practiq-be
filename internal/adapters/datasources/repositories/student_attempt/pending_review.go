package studentattempt

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// PendingReviewFilter narrows the queue. Empty strings mean "no filter", the
// same convention the notebook submission queue uses.
type PendingReviewFilter struct {
	TeacherID string
	CourseID  string
	StudentID string
	// SheetType narrows the queue further, but the queue only ever contains
	// level tests: a practice is graded on submit and never handed to a
	// teacher, so filtering by "practice" returns nothing.
	SheetType string
	// Reviewed is "", "reviewed" or "unreviewed".
	Reviewed string
	Limit    int
	Offset   int
}

// ListPendingReview returns answers requiring a teacher from that teacher's
// courses: level test answers the assistant could not resolve. It used to admit
// every answer carrying an attachment, which put files the assistant had
// already graded in front of the teacher. Practice answers never qualify — a
// practice is graded on submit and must not wait on a correction — and homework
// has its own queue (see the notebook submission repository).
func (r *repository) ListPendingReview(ctx context.Context, filter PendingReviewFilter) ([]domain.PendingAttemptReview, error) {
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
		WHERE sa.needs_teacher_review
		  AND COALESCE(ps.sheet_type, '') = 'level_test'
		  AND c.teacher_id = $1
		  AND c.deleted_at IS NULL
		  AND ($2 = '' OR c.id::text = $2)
		  AND ($3 = '' OR sa.student_id = $3)
		  AND ($4 = '' OR COALESCE(ps.sheet_type, '') = $4)
		  AND ($5 = '' OR ($5 = 'reviewed' AND sa.teacher_reviewed_at IS NOT NULL)
		               OR ($5 = 'unreviewed' AND sa.teacher_reviewed_at IS NULL))
		ORDER BY sa.created_at DESC
	`
	args := []any{filter.TeacherID, filter.CourseID, filter.StudentID, filter.SheetType, filter.Reviewed}

	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	// One row past the page so the caller can tell whether another exists; a
	// full page is not proof of it.
	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit+1, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
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
