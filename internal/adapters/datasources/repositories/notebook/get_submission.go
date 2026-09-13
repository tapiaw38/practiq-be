package notebook

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetSubmission(ctx context.Context, pageID, studentID string) (*domain.NotebookSubmission, error) {
	var s domain.NotebookSubmission
	err := r.db.QueryRowContext(ctx, `
		SELECT ns.id, ns.page_id, ns.student_id, ns.canvas_data, ns.answer_text, COALESCE(ns.ai_recognized_text,''), ns.ai_is_correct, ns.ai_feedback, ns.ai_reviewed_at, COALESCE(ns.needs_teacher_review,FALSE),
		       ns.teacher_is_correct, COALESCE(ns.teacher_feedback,''), ns.teacher_reviewed_at, ns.submitted_at, ns.updated_at
		FROM notebook_submissions ns
		JOIN notebook_pages np ON np.id = ns.page_id
		JOIN notebooks n ON n.id = np.notebook_id
		JOIN courses c ON c.id = n.course_id
		WHERE ns.page_id = $1 AND ns.student_id = $2 AND n.deleted_at IS NULL AND c.deleted_at IS NULL
	`, pageID, studentID).Scan(&s.ID, &s.PageID, &s.StudentID, &s.CanvasData, &s.AnswerText, &s.AIRecognizedText, &s.AIIsCorrect, &s.AIFeedback, &s.AIReviewedAt, &s.NeedsTeacherReview, &s.TeacherIsCorrect, &s.TeacherFeedback, &s.TeacherReviewedAt, &s.SubmittedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}
