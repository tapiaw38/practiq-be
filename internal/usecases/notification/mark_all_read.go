package notification

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	MarkAllReadUsecase interface {
		Execute(ctx context.Context, userID string) (*MarkReadOutput, apperrors.ApplicationError)
	}

	markAllReadUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewMarkAllReadUsecase(contextFactory appcontext.Factory) MarkAllReadUsecase {
	return &markAllReadUsecase{contextFactory: contextFactory}
}

func (u *markAllReadUsecase) Execute(ctx context.Context, userID string) (*MarkReadOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if err := app.Repositories.Notification.MarkAllRead(ctx, userID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotificationUpdateError, err)
	}
	return &MarkReadOutput{Data: OperationResultData{Message: "notifications marked as read"}}, nil
}
