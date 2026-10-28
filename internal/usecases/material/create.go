package material

import (
	"context"
	"errors"
	"strings"

	materialRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/material"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID, teacherID string, in CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Title         string `json:"title" binding:"required"`
		Type          string `json:"type" binding:"required"`
		ExtractedText string `json:"extracted_text"`
		FileURL       string `json:"file_url"`
	}

	CreateOutput struct {
		Data MaterialData `json:"data"`
	}
)

func materialStatus(fileURL string) string {
	if strings.TrimSpace(fileURL) == "" {
		return "text_only"
	}
	return "uploaded"
}

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID, teacherID string, in CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := requesterCanWriteCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
	}

	if in.FileURL != "" && (app.ImageStorage == nil ||
		!app.ImageStorage.OwnsFileURL(in.FileURL, materialsFolder, teacherID)) {
		return nil, apperrors.NewApplicationError(mappings.MaterialCreateError,
			errors.New("the file does not belong to this teacher"))
	}

	id, err := app.Repositories.Material.Create(ctx, domain.Material{
		CourseID:      courseID,
		TeacherID:     teacherID,
		Title:         in.Title,
		Type:          in.Type,
		ExtractedText: in.ExtractedText,
		FileURL:       in.FileURL,
		Status:        materialStatus(in.FileURL),
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialCreateError, err)
	}

	materials, err := app.Repositories.Material.List(ctx, materialRepo.ListFilter{CourseID: courseID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.MaterialListError, err)
	}
	for _, m := range materials {
		if m.ID == id {
			return &CreateOutput{Data: withViewURL(app, toMaterialData(m))}, nil
		}
	}

	return nil, apperrors.NewInternalError(nil)
}
