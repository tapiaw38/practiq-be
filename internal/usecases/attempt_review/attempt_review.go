package attemptreview

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

const timeFormat = "2006-01-02T15:04:05Z"

type (
	ListUsecase interface {
		Execute(ctx context.Context, teacherID string, includeReviewed bool) (*ListOutput, apperrors.ApplicationError)
	}

	ReviewUsecase interface {
		Execute(ctx context.Context, attemptID, teacherID string, isAdmin bool, input ReviewInput) (*ReviewOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	reviewUsecase struct {
		contextFactory appcontext.Factory
	}

	ReviewInput struct {
		IsCorrect bool
		Feedback  string
	}

	ReviewData struct {
		AttemptID             string `json:"attempt_id"`
		StudentID             string `json:"student_id"`
		StudentName           string `json:"student_name,omitempty"`
		ExerciseID            string `json:"exercise_id"`
		Question              string `json:"question"`
		ExerciseType          string `json:"exercise_type"`
		PracticeSheetID       string `json:"practice_sheet_id,omitempty"`
		PracticeSheetTitle    string `json:"practice_sheet_title,omitempty"`
		SheetType             string `json:"sheet_type,omitempty"`
		CourseID              string `json:"course_id"`
		CourseTitle           string `json:"course_title"`
		AttachmentURL         string `json:"attachment_url,omitempty"`
		AttachmentName        string `json:"attachment_name,omitempty"`
		AttachmentContentType string `json:"attachment_content_type,omitempty"`
		AnswerText            string `json:"answer_text,omitempty"`
		AIFeedback            string `json:"ai_feedback,omitempty"`
		TeacherIsCorrect      *bool  `json:"teacher_is_correct,omitempty"`
		TeacherFeedback       string `json:"teacher_feedback,omitempty"`
		TeacherReviewedAt     string `json:"teacher_reviewed_at,omitempty"`
		CreatedAt             string `json:"created_at"`
	}

	ListOutput struct {
		Data []ReviewData `json:"data"`
	}

	ReviewOutput struct {
		Data OperationResultData `json:"data"`
	}

	OperationResultData struct {
		Message string `json:"message"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func NewReviewUsecase(contextFactory appcontext.Factory) ReviewUsecase {
	return &reviewUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, teacherID string, includeReviewed bool) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	reviews, err := app.Repositories.StudentAttempt.ListPendingReview(ctx, teacherID, includeReviewed)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewListError, err)
	}

	data := make([]ReviewData, 0, len(reviews))
	for _, review := range reviews {
		data = append(data, toReviewData(review))
	}
	return &ListOutput{Data: data}, nil
}

func (u *reviewUsecase) Execute(ctx context.Context, attemptID, teacherID string, isAdmin bool, input ReviewInput) (*ReviewOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	owner, err := app.Repositories.StudentAttempt.GetTeacherForAttempt(ctx, attemptID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewError, err)
	}
	if owner == "" {
		return nil, apperrors.NewNotFoundError("attempt not found")
	}
	if !isAdmin && owner != teacherID {
		return nil, apperrors.NewForbiddenError()
	}

	if err := app.Repositories.StudentAttempt.Review(ctx, attemptID, input.IsCorrect, input.Feedback); err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewError, err)
	}
	return &ReviewOutput{Data: OperationResultData{Message: "attempt reviewed"}}, nil
}

func toReviewData(r domain.PendingAttemptReview) ReviewData {
	data := ReviewData{
		AttemptID:             r.AttemptID,
		StudentID:             r.StudentID,
		StudentName:           r.StudentName,
		ExerciseID:            r.ExerciseID,
		Question:              r.Question,
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
		TeacherIsCorrect:      r.TeacherIsCorrect,
		TeacherFeedback:       r.TeacherFeedback,
		CreatedAt:             r.CreatedAt.UTC().Format(timeFormat),
	}
	if r.TeacherReviewedAt != nil {
		data.TeacherReviewedAt = r.TeacherReviewedAt.UTC().Format(timeFormat)
	}
	return data
}
