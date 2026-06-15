package notebook

import (
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type (
	SubmissionData struct {
		ID                string `json:"id"`
		CanvasData        string `json:"canvas_data"`
		AnswerText        string `json:"answer_text"`
		AIRecognizedText  string `json:"ai_recognized_text,omitempty"`
		AIIsCorrect       *bool  `json:"ai_is_correct,omitempty"`
		AIFeedback        string `json:"ai_feedback,omitempty"`
		AIReviewedAt      string `json:"ai_reviewed_at,omitempty"`
		TeacherIsCorrect  *bool  `json:"teacher_is_correct,omitempty"`
		TeacherFeedback   string `json:"teacher_feedback,omitempty"`
		TeacherReviewedAt string `json:"teacher_reviewed_at,omitempty"`
	}

	NotebookSubmissionFullData struct {
		ID                string `json:"id"`
		PageID            string `json:"page_id"`
		StudentID         string `json:"student_id"`
		StudentName       string `json:"student_name,omitempty"`
		StudentEmail      string `json:"student_email,omitempty"`
		NotebookID        string `json:"notebook_id"`
		NotebookTitle     string `json:"notebook_title,omitempty"`
		PageTitle         string `json:"page_title,omitempty"`
		PageNumber        int    `json:"page_number"`
		CourseID          string `json:"course_id"`
		CanvasData        string `json:"canvas_data"`
		AnswerText        string `json:"answer_text"`
		AIRecognizedText  string `json:"ai_recognized_text,omitempty"`
		AIIsCorrect       *bool  `json:"ai_is_correct,omitempty"`
		AIFeedback        string `json:"ai_feedback,omitempty"`
		AIReviewedAt      string `json:"ai_reviewed_at,omitempty"`
		TeacherIsCorrect  *bool  `json:"teacher_is_correct,omitempty"`
		TeacherFeedback   string `json:"teacher_feedback,omitempty"`
		TeacherReviewedAt string `json:"teacher_reviewed_at,omitempty"`
		CreatedAt         string `json:"created_at"`
		UpdatedAt         string `json:"updated_at,omitempty"`
	}

	PageData struct {
		ID           string          `json:"id"`
		NotebookID   string          `json:"notebook_id"`
		PageNumber   int             `json:"page_number"`
		Title        string          `json:"title"`
		ContentType  string          `json:"content_type"`
		ContentData  string          `json:"content_data"`
		Instructions string          `json:"instructions"`
		Submission   *SubmissionData `json:"submission,omitempty"`
	}

	NotebookData struct {
		ID          string     `json:"id"`
		CourseID    string     `json:"course_id"`
		TeacherID   string     `json:"teacher_id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Level       int        `json:"level"`
		Pages       []PageData `json:"pages"`
		CreatedAt   string     `json:"created_at"`
	}
)

func toNotebookData(nb *domain.Notebook) NotebookData {
	pages := make([]PageData, 0, len(nb.Pages))
	for _, p := range nb.Pages {
		po := PageData{
			ID:           p.ID,
			NotebookID:   p.NotebookID,
			PageNumber:   p.PageNumber,
			Title:        p.Title,
			ContentType:  p.ContentType,
			ContentData:  p.ContentData,
			Instructions: p.Instructions,
		}
		if p.Submission != nil {
			aiReviewedAt := ""
			if p.Submission.AIReviewedAt != nil {
				aiReviewedAt = p.Submission.AIReviewedAt.Format("2006-01-02T15:04:05Z")
			}
			po.Submission = &SubmissionData{
				ID:                p.Submission.ID,
				CanvasData:        p.Submission.CanvasData,
				AnswerText:        p.Submission.AnswerText,
				AIRecognizedText:  p.Submission.AIRecognizedText,
				AIIsCorrect:       p.Submission.AIIsCorrect,
				AIFeedback:        p.Submission.AIFeedback,
				AIReviewedAt:      aiReviewedAt,
				TeacherIsCorrect:  p.Submission.TeacherIsCorrect,
				TeacherFeedback:   p.Submission.TeacherFeedback,
				TeacherReviewedAt: formatTimePtr(p.Submission.TeacherReviewedAt),
			}
		}
		pages = append(pages, po)
	}
	return NotebookData{
		ID:          nb.ID,
		CourseID:    nb.CourseID,
		TeacherID:   nb.TeacherID,
		Title:       nb.Title,
		Description: nb.Description,
		Level:       nb.Level,
		Pages:       pages,
		CreatedAt:   nb.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toPageData(id string, input AddPageInput, contentData string) PageData {
	return PageData{
		ID:           id,
		NotebookID:   input.NotebookID,
		PageNumber:   input.PageNumber,
		Title:        input.Title,
		ContentType:  input.ContentType,
		ContentData:  contentData,
		Instructions: input.Instructions,
	}
}

func toFullSubmissionData(s domain.NotebookSubmissionFull) NotebookSubmissionFullData {
	return NotebookSubmissionFullData{
		ID:                s.ID,
		PageID:            s.PageID,
		StudentID:         s.StudentID,
		StudentName:       s.StudentName,
		StudentEmail:      s.StudentEmail,
		NotebookID:        s.NotebookID,
		NotebookTitle:     s.NotebookTitle,
		PageTitle:         s.PageTitle,
		PageNumber:        s.PageNumber,
		CourseID:          s.CourseID,
		CanvasData:        s.CanvasData,
		AnswerText:        s.AnswerText,
		AIRecognizedText:  s.AIRecognizedText,
		AIIsCorrect:       s.AIIsCorrect,
		AIFeedback:        s.AIFeedback,
		AIReviewedAt:      formatTimePtr(s.AIReviewedAt),
		TeacherIsCorrect:  s.TeacherIsCorrect,
		TeacherFeedback:   s.TeacherFeedback,
		TeacherReviewedAt: formatTimePtr(s.TeacherReviewedAt),
		CreatedAt:         s.SubmittedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:         s.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func formatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z")
}
