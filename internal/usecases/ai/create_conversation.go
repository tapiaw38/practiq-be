package ai

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type CreateConversationUsecase interface {
	Execute(context.Context, CreateConversationInput) (*ConversationOutput, apperrors.ApplicationError)
}

type createConversationUsecase struct {
	factory appcontext.Factory
}

type CreateConversationInput struct {
	StudentID       string
	CourseID        string `json:"course_id"`
	PracticeSheetID string `json:"practice_sheet_id"`
}

func NewCreateConversationUsecase(factory appcontext.Factory) CreateConversationUsecase {
	return &createConversationUsecase{factory: factory}
}

func (u *createConversationUsecase) Execute(ctx context.Context, input CreateConversationInput) (*ConversationOutput, apperrors.ApplicationError) {
	app := u.factory()

	id, err := app.Repositories.AIConversation.CreateConversation(ctx, domain.AIConversation{
		StudentID:       input.StudentID,
		CourseID:        input.CourseID,
		PracticeSheetID: input.PracticeSheetID,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AIConversationCreateError, err)
	}

	return &ConversationOutput{Data: ConversationData{
		ID:              id,
		StudentID:       input.StudentID,
		CourseID:        input.CourseID,
		PracticeSheetID: input.PracticeSheetID,
	}}, nil
}
