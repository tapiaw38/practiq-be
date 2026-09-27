package attemptreview

import (
	"context"
	"time"

	studentAttemptRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/student_attempt"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

const timeFormat = "2006-01-02T15:04:05Z"

const attachmentLinkTTL = time.Hour

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

type (
	ListUsecase interface {
		Execute(ctx context.Context, teacherID string, input ListInput) (*ListOutput, apperrors.ApplicationError)
	}

	ReviewUsecase interface {
		Execute(ctx context.Context, attemptID, teacherID string, isSuperAdmin bool, input ReviewInput) (*ReviewOutput, apperrors.ApplicationError)
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

	ListInput struct {
		CourseID    string
		StudentID   string
		SheetType   string
		Reviewed    string
		Limit       int
		Offset      int
		BearerToken string
	}

	ListOutput struct {
		Data []ReviewData `json:"data"`

		HasMore bool `json:"has_more"`
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

func (u *listUsecase) Execute(ctx context.Context, teacherID string, input ListInput) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	limit := input.Limit
	if limit <= 0 || limit > maxPageSize {
		limit = defaultPageSize
	}

	reviews, err := app.Repositories.StudentAttempt.ListPendingReview(ctx, studentAttemptRepo.PendingReviewFilter{
		TeacherID: teacherID,
		CourseID:  input.CourseID,
		StudentID: input.StudentID,
		SheetType: input.SheetType,
		Reviewed:  input.Reviewed,
		Limit:     limit,
		Offset:    input.Offset,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewListError, err)
	}

	hasMore := len(reviews) > limit
	if hasMore {
		reviews = reviews[:limit]
	}

	ids := make([]string, 0, len(reviews))
	for _, review := range reviews {
		ids = append(ids, review.StudentID)
	}
	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, ids)
	if appErr != nil {
		return nil, appErr
	}
	for i, review := range reviews {
		reviews[i].StudentName = identity.FullName(names[review.StudentID], review.StudentID)
	}

	data := make([]ReviewData, 0, len(reviews))
	for _, review := range reviews {
		item := toReviewData(review)
		if item.AttachmentURL != "" && app.ImageStorage != nil {
			if signed, ok := app.ImageStorage.PresignGetURL(item.AttachmentURL, attachmentLinkTTL); ok {
				item.AttachmentViewURL = signed
			}
		}
		if review.ImageURL != "" && app.ImageStorage != nil {
			if signed, ok := app.ImageStorage.PresignGetURL(review.ImageURL, attachmentLinkTTL); ok {
				item.ImageViewURL = signed
			}
		}
		if review.StatementMediaURL != "" && app.ImageStorage != nil {
			if signed, ok := app.ImageStorage.PresignGetURL(review.StatementMediaURL, attachmentLinkTTL); ok {
				item.StatementMediaViewURL = signed
			}
		}
		data = append(data, item)
	}
	return &ListOutput{Data: data, HasMore: hasMore}, nil
}

func (u *reviewUsecase) Execute(ctx context.Context, attemptID, teacherID string, isSuperAdmin bool, input ReviewInput) (*ReviewOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	owner, err := app.Repositories.StudentAttempt.GetTeacherForAttempt(ctx, attemptID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewError, err)
	}
	if owner == "" {
		return nil, apperrors.NewNotFoundError("attempt not found")
	}
	if !isSuperAdmin && owner != teacherID {
		return nil, apperrors.NewForbiddenError()
	}

	attemptCtx, err := app.Repositories.StudentAttempt.GetAttemptContext(ctx, attemptID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewError, err)
	}
	if attemptCtx.SheetType != sheetTypeLevelTest {
		return nil, apperrors.NewBadRequestError("only level test answers are corrected by the teacher")
	}

	if err := app.Repositories.StudentAttempt.Review(ctx, attemptID, input.IsCorrect, input.Feedback); err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewError, err)
	}

	applyLevelTestOutcome(ctx, app, attemptID)

	return &ReviewOutput{Data: OperationResultData{Message: "attempt reviewed"}}, nil
}

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
