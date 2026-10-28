package attemptreview

import (
	"context"

	studentAttemptRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/student_attempt"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

type (
	ListUsecase interface {
		Execute(ctx context.Context, teacherID string, input ListInput) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
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
		Data    []ReviewData `json:"data"`
		HasMore bool         `json:"has_more"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
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
