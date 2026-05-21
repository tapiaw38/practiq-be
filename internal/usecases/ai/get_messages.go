package ai

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GetMessagesUsecase interface {
	Execute(context.Context, string) (*MessagesListOutput, apperrors.ApplicationError)
}

type getMessagesUsecase struct {
	factory appcontext.Factory
}

func NewGetMessagesUsecase(factory appcontext.Factory) GetMessagesUsecase {
	return &getMessagesUsecase{factory: factory}
}

func (u *getMessagesUsecase) Execute(ctx context.Context, conversationID string) (*MessagesListOutput, apperrors.ApplicationError) {
	app := u.factory()

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

	return &MessagesListOutput{Data: data}, nil
}
