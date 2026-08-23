package notebook

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetFullSubmissionByID(ctx context.Context, id string) (*domain.NotebookSubmissionFull, error) {
	var s domain.NotebookSubmissionFull
	err := r.db.QueryRowContext(ctx, `
		SELECT ns.id, ns.page_id, ns.student_id, ns.canvas_data, ns.answer_text,
		       COALESCE(ns.ai_recognized_text,''), ns.ai_is_correct, COALESCE(ns.ai_feedback,''), ns.ai_reviewed_at, COALESCE(ns.needs_teacher_review,FALSE),
		       ns.teacher_is_correct, COALESCE(ns.teacher_feedback,''), ns.teacher_reviewed_at,
		       ns.submitted_at, ns.updated_at,
		       COALESCE(up.name,''), COALESCE(up.email,''),
		       n.id, n.title, COALESCE(np.title,''), np.page_number, n.course_id::text, n.teacher_id
		FROM notebook_submissions ns
		JOIN notebook_pages np ON np.id = ns.page_id
		JOIN notebooks n ON n.id = np.notebook_id
		JOIN courses c ON c.id = n.course_id
		LEFT JOIN user_profiles up ON up.id = ns.student_id
		WHERE ns.id = $1 AND n.deleted_at IS NULL AND c.deleted_at IS NULL
	`, id).Scan(
		&s.ID, &s.PageID, &s.StudentID, &s.CanvasData, &s.AnswerText,
		&s.AIRecognizedText, &s.AIIsCorrect, &s.AIFeedback, &s.AIReviewedAt, &s.NeedsTeacherReview,
		&s.TeacherIsCorrect, &s.TeacherFeedback, &s.TeacherReviewedAt,
		&s.SubmittedAt, &s.UpdatedAt,
		&s.StudentName, &s.StudentEmail,
		&s.NotebookID, &s.NotebookTitle, &s.PageTitle, &s.PageNumber, &s.CourseID, &s.TeacherID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}
