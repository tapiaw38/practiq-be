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
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		GradeID     string `json:"grade_id"`
		SubjectID   string `json:"subject_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Level       string `json:"level"`
		Subject     string `json:"subject"`
		Status      string `json:"status"`
	}

	UpdateOutput struct {
		Data CourseData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if strings.TrimSpace(in.GradeID) == "" {
		return nil, apperrors.NewBadRequestError("grade_id is required")
	}
	if strings.TrimSpace(in.SubjectID) == "" {
		return nil, apperrors.NewBadRequestError("subject_id is required")
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

	status := strings.TrimSpace(in.Status)
	if status == "" {
		status = current.Status
	}
	switch status {
	case domain.CourseStatusDraft, domain.CourseStatusPublished, domain.CourseStatusArchived:
	default:
		return nil, apperrors.NewBadRequestError("status must be draft, published or archived")
	}

	if err := app.Repositories.Course.Update(ctx, id, domain.Course{
		GradeID:     in.GradeID,
		SubjectID:   in.SubjectID,
		Title:       in.Title,
		Description: in.Description,
		Level:       in.Level,
		Subject:     in.Subject,
		Status:      status,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseUpdateError, err)
	}

	c, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}

	return &UpdateOutput{Data: toCourseData(*c)}, nil
}
