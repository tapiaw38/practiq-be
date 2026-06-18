package practicesheet

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Title       string
		TopicID     string
		Level       int
		SheetType   string
		TestStyle   string
		ExerciseIDs []string
	}

	UpdateOutput struct {
		Data PracticeSheetData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Verify practice sheet exists and check ownership
	ps, err := app.Repositories.PracticeSheet.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	// Check course ownership for authorization
	if !isAdmin {
		course, err := app.Repositories.Course.Get(ctx, ps.CourseID)
		if err != nil || course == nil {
			return nil, apperrors.NewForbiddenError()
		}
		if course.TeacherID != requesterID {
			return nil, apperrors.NewForbiddenError()
		}
	}

	if err := app.Repositories.PracticeSheet.Update(ctx, id, domain.PracticeSheet{
		Title:     input.Title,
		TopicID:   input.TopicID,
		Level:     input.Level,
		SheetType: input.SheetType,
		TestStyle: input.TestStyle,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetUpdateError, err)
	}

	if err := app.Repositories.PracticeSheet.ReplaceExercises(ctx, id, input.ExerciseIDs); err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetUpdateError, err)
	}

	ps, err = app.Repositories.PracticeSheet.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	return &UpdateOutput{Data: toSheetData(*ps)}, nil
}
