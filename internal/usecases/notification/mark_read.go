package notification

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	MarkReadUsecase interface {
		Execute(ctx context.Context, id, userID string) (*MarkReadOutput, apperrors.ApplicationError)
	}

	markReadUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewMarkReadUsecase(contextFactory appcontext.Factory) MarkReadUsecase {
	return &markReadUsecase{contextFactory: contextFactory}
}

func (u *markReadUsecase) Execute(ctx context.Context, id, userID string) (*MarkReadOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	updated, err := app.Repositories.Notification.MarkRead(ctx, id, userID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotificationUpdateError, err)
	}
	if !updated {
		return &MarkReadOutput{Data: OperationResultData{Message: "notification already read"}}, nil
	}
	return &MarkReadOutput{Data: OperationResultData{Message: "notification marked as read"}}, nil
}
