package learningstrategy

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
	AssignToCourseUsecase interface {
		Execute(ctx context.Context, requesterID, courseID string, isSuperAdmin bool, input AssignToCourseInput) (*AssignToCourseOutput, apperrors.ApplicationError)
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

func (u *assignToCourseUsecase) Execute(ctx context.Context, requesterID, courseID string, isSuperAdmin bool, input AssignToCourseInput) (*AssignToCourseOutput, apperrors.ApplicationError) {
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

	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
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
