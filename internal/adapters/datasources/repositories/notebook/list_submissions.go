package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListSubmissions(ctx context.Context, filter SubmissionFilter) ([]domain.NotebookSubmissionFull, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ns.id, ns.page_id, ns.student_id, ns.canvas_data, ns.answer_text,
		       COALESCE(ns.ai_recognized_text,''), ns.ai_is_correct, COALESCE(ns.ai_feedback,''), ns.ai_reviewed_at,
		       ns.teacher_is_correct, COALESCE(ns.teacher_feedback,''), ns.teacher_reviewed_at,
		       ns.submitted_at, ns.updated_at,
		       COALESCE(up.name,''), COALESCE(up.email,''),
		       n.id, n.title, COALESCE(np.title,''), np.page_number, n.course_id::text, n.teacher_id
		FROM notebook_submissions ns
		JOIN notebook_pages np ON np.id = ns.page_id
		JOIN notebooks n ON n.id = np.notebook_id
		LEFT JOIN user_profiles up ON up.id = ns.student_id
		WHERE ($1 = '' OR n.id::text = $1)
		  AND ($2 = '' OR ns.student_id = $2)
		  AND ($3 = '' OR n.course_id::text = $3)
		  AND ($4 = '' OR ($4 = 'reviewed' AND ns.ai_reviewed_at IS NOT NULL) OR ($4 = 'unreviewed' AND ns.ai_reviewed_at IS NULL))
		  AND ($5 = '' OR n.teacher_id = $5)
		ORDER BY ns.submitted_at DESC
	`, filter.NotebookID, filter.StudentID, filter.CourseID, filter.Reviewed, filter.TeacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.NotebookSubmissionFull
	for rows.Next() {
		var s domain.NotebookSubmissionFull
		if err := rows.Scan(
			&s.ID, &s.PageID, &s.StudentID, &s.CanvasData, &s.AnswerText,
			&s.AIRecognizedText, &s.AIIsCorrect, &s.AIFeedback, &s.AIReviewedAt,
			&s.TeacherIsCorrect, &s.TeacherFeedback, &s.TeacherReviewedAt,
			&s.SubmittedAt, &s.UpdatedAt,
			&s.StudentName, &s.StudentEmail,
			&s.NotebookID, &s.NotebookTitle, &s.PageTitle, &s.PageNumber, &s.CourseID, &s.TeacherID,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}
