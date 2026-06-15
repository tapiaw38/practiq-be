package notebook

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) GetSubmission(ctx context.Context, pageID, studentID string) (*domain.NotebookSubmission, error) {
	var s domain.NotebookSubmission
	err := r.db.QueryRowContext(ctx, `
		SELECT id, page_id, student_id, canvas_data, answer_text, COALESCE(ai_recognized_text,''), ai_is_correct, ai_feedback, ai_reviewed_at,
		       teacher_is_correct, COALESCE(teacher_feedback,''), teacher_reviewed_at, submitted_at, updated_at
		FROM notebook_submissions WHERE page_id = $1 AND student_id = $2
	`, pageID, studentID).Scan(&s.ID, &s.PageID, &s.StudentID, &s.CanvasData, &s.AnswerText, &s.AIRecognizedText, &s.AIIsCorrect, &s.AIFeedback, &s.AIReviewedAt, &s.TeacherIsCorrect, &s.TeacherFeedback, &s.TeacherReviewedAt, &s.SubmittedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}
