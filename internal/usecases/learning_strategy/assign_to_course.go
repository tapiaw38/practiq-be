package learningstrategy

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	AssignToCourseUsecase interface {
		Execute(ctx context.Context, requesterID, courseID string, isAdmin bool, input AssignToCourseInput) (*AssignToCourseOutput, apperrors.ApplicationError)
	}

	assignToCourseUsecase struct {
		contextFactory appcontext.Factory
	}

	AssignToCourseInput struct {
		StrategyID string `json:"strategy_id"`
		IsDefault  bool   `json:"is_default"`
		Config     string `json:"config"`
	}

	AssignToCourseOutput struct {
		Data CourseLearningStrategyData `json:"data"`
	}
)

func NewAssignToCourseUsecase(contextFactory appcontext.Factory) AssignToCourseUsecase {
	return &assignToCourseUsecase{contextFactory: contextFactory}
}

func (u *assignToCourseUsecase) Execute(ctx context.Context, requesterID, courseID string, isAdmin bool, input AssignToCourseInput) (*AssignToCourseOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if strings.TrimSpace(input.StrategyID) == "" {
		return nil, apperrors.NewBadRequestError("strategy_id is required")
	}

	// Verify strategy exists
	strategy, err := app.Repositories.LearningStrategy.Get(ctx, input.StrategyID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}
	if strategy == nil {
		return nil, apperrors.NewNotFoundError("learning strategy not found")
	}

	// Verify course exists and check ownership
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return nil, apperrors.NewNotFoundError("course not found")
	}

	if !isAdmin && course.TeacherID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}

	config := input.Config
	if config == "" {
		config = "{}"
	}

	id, err := app.Repositories.LearningStrategy.AssignToCourse(ctx, domain.CourseLearningStrategy{
		CourseID:   courseID,
		StrategyID: input.StrategyID,
		IsDefault:  input.IsDefault,
		Config:     config,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyAssignError, err)
	}

	cls, err := app.Repositories.LearningStrategy.GetCourseStrategy(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}

	return &AssignToCourseOutput{Data: toCourseStrategyData(*cls)}, nil
}
