package ai

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	CreateConversationUsecase interface {
		Execute(ctx context.Context, studentID string, in CreateConversationInput) (*CreateConversationOutput, apperrors.ApplicationError)
	}

	createConversationUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateConversationInput struct {
		CourseID        string `json:"course_id"`
		PracticeSheetID string `json:"practice_sheet_id"`
	}

	CreateConversationOutput struct {
		Data ConversationData `json:"data"`
	}
)

func NewCreateConversationUsecase(contextFactory appcontext.Factory) CreateConversationUsecase {
	return &createConversationUsecase{contextFactory: contextFactory}
}

func (u *createConversationUsecase) Execute(ctx context.Context, studentID string, in CreateConversationInput) (*CreateConversationOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	id, err := app.Repositories.AIConversation.CreateConversation(ctx, domain.AIConversation{
		StudentID:       studentID,
		CourseID:        in.CourseID,
		PracticeSheetID: in.PracticeSheetID,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AIConversationCreateError, err)
	}

	return &CreateConversationOutput{Data: toConversationOutputData(id, studentID, in)}, nil
}
