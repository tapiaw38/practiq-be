package course

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	SetStatusUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id, status string) (*UpdateOutput, apperrors.ApplicationError)
	}

	setStatusUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewSetStatusUsecase(contextFactory appcontext.Factory) SetStatusUsecase {
	return &setStatusUsecase{contextFactory: contextFactory}
}

// ValidCourseStatus reports whether a status is one a course can hold.
func ValidCourseStatus(status string) bool {
	switch status {
	case domain.CourseStatusDraft, domain.CourseStatusPublished, domain.CourseStatusArchived:
		return true
	default:
		return false
	}
}

// Execute moves a course through its lifecycle and touches nothing else.
//
// Update demands the whole course — grade, subject, title — because it writes
// every column. Publishing is not an edit of those: asking the teacher's
// browser to resend a course it is not changing is how a lifecycle change
// blanks a description.
func (u *setStatusUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id, status string) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	status = strings.TrimSpace(status)
	if !ValidCourseStatus(status) {
		return nil, apperrors.NewBadRequestError("status must be draft, published or archived")
	}

	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, id); appErr != nil {
		return nil, appErr
	}

	current, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if current == nil {
		return nil, apperrors.NewNotFoundError("course not found")
	}

	if current.Status != status {
		current.Status = status
		if err := app.Repositories.Course.Update(ctx, id, *current); err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseUpdateError, err)
		}
	}

	updated, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}

	return &UpdateOutput{Data: toCourseData(*updated)}, nil
}
