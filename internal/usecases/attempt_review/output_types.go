package attemptreview

import (
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

const timeFormat = "2006-01-02T15:04:05Z"

const attachmentLinkTTL = time.Hour

type (
	ReviewData struct {
		AttemptID             string `json:"attempt_id"`
		StudentID             string `json:"student_id"`
		StudentName           string `json:"student_name,omitempty"`
		ExerciseID            string `json:"exercise_id"`
		Question              string `json:"question"`
		ExerciseType          string `json:"exercise_type"`
		StatementMediaViewURL string `json:"statement_media_view_url,omitempty"`

		HasTeacherImage    bool   `json:"has_teacher_image,omitempty"`
		PracticeSheetID    string `json:"practice_sheet_id,omitempty"`
		PracticeSheetTitle string `json:"practice_sheet_title,omitempty"`
		SheetType          string `json:"sheet_type,omitempty"`
		CourseID           string `json:"course_id"`
		CourseTitle        string `json:"course_title"`

		ImageViewURL  string `json:"image_view_url,omitempty"`
		AttachmentURL string `json:"attachment_url,omitempty"`

		AttachmentViewURL     string `json:"attachment_view_url,omitempty"`
		AttachmentName        string `json:"attachment_name,omitempty"`
		AttachmentContentType string `json:"attachment_content_type,omitempty"`
		AnswerText            string `json:"answer_text,omitempty"`
		AIFeedback            string `json:"ai_feedback,omitempty"`

		AIIsCorrect       *bool  `json:"ai_is_correct,omitempty"`
		TeacherIsCorrect  *bool  `json:"teacher_is_correct,omitempty"`
		TeacherFeedback   string `json:"teacher_feedback,omitempty"`
		TeacherReviewedAt string `json:"teacher_reviewed_at,omitempty"`
		CreatedAt         string `json:"created_at"`
	}

	OperationResultData struct {
		Message string `json:"message"`
	}
)

func toReviewData(r domain.PendingAttemptReview) ReviewData {
	data := ReviewData{
		AttemptID:             r.AttemptID,
		StudentID:             r.StudentID,
		StudentName:           r.StudentName,
		ExerciseID:            r.ExerciseID,
		Question:              r.Question,
		HasTeacherImage:       r.HasTeacherImage,
		ExerciseType:          r.ExerciseType,
		PracticeSheetID:       r.PracticeSheetID,
		PracticeSheetTitle:    r.PracticeSheetTitle,
		SheetType:             r.SheetType,
		CourseID:              r.CourseID,
		CourseTitle:           r.CourseTitle,
		AttachmentURL:         r.AttachmentURL,
		AttachmentName:        r.AttachmentName,
		AttachmentContentType: r.AttachmentContentType,
		AnswerText:            r.AnswerText,
		AIFeedback:            r.AIFeedback,
		AIIsCorrect:           r.AIIsCorrect,
		TeacherIsCorrect:      r.TeacherIsCorrect,
		TeacherFeedback:       r.TeacherFeedback,
		CreatedAt:             r.CreatedAt.UTC().Format(timeFormat),
	}
	if r.TeacherReviewedAt != nil {
		data.TeacherReviewedAt = r.TeacherReviewedAt.UTC().Format(timeFormat)
	}
	return data
}
