package ai

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetMessagesUsecase interface {
		Execute(context.Context, string) (*GetMessagesOutput, apperrors.ApplicationError)
	}

	getMessagesUsecase struct {
		contextFactory appcontext.Factory
	}

	GetMessagesOutput struct {
		Data []MessageData `json:"data"`
	}
)

func NewGetMessagesUsecase(contextFactory appcontext.Factory) GetMessagesUsecase {
	return &getMessagesUsecase{contextFactory: contextFactory}
}

func (u *getMessagesUsecase) Execute(ctx context.Context, conversationID string) (*GetMessagesOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	messages, err := app.Repositories.AIConversation.ListMessages(ctx, conversationID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AIMessageListError, err)
	}

	var data []MessageData
	for _, m := range messages {
		data = append(data, toMessageData(m))
	}
	if data == nil {
		data = []MessageData{}
	}

	return &GetMessagesOutput{Data: data}, nil
}
