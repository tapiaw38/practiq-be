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

	MarkAllReadUsecase interface {
		Execute(ctx context.Context, userID string) (*MarkReadOutput, apperrors.ApplicationError)
	}

	markReadUsecase struct {
		contextFactory appcontext.Factory
	}

	markAllReadUsecase struct {
		contextFactory appcontext.Factory
	}

	MarkReadOutput struct {
		Data OperationResultData `json:"data"`
	}
)

func NewMarkReadUsecase(contextFactory appcontext.Factory) MarkReadUsecase {
	return &markReadUsecase{contextFactory: contextFactory}
}

func NewMarkAllReadUsecase(contextFactory appcontext.Factory) MarkAllReadUsecase {
	return &markAllReadUsecase{contextFactory: contextFactory}
}

// Execute marks one notification as read. Scoped by user, so a notification
// that belongs to someone else is indistinguishable from a missing one.
func (u *markReadUsecase) Execute(ctx context.Context, id, userID string) (*MarkReadOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	updated, err := app.Repositories.Notification.MarkRead(ctx, id, userID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotificationUpdateError, err)
	}
	if !updated {
		// Already read is not an error: the client may retry or double-click.
		return &MarkReadOutput{Data: OperationResultData{Message: "notification already read"}}, nil
	}
	return &MarkReadOutput{Data: OperationResultData{Message: "notification marked as read"}}, nil
}

func (u *markAllReadUsecase) Execute(ctx context.Context, userID string) (*MarkReadOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if err := app.Repositories.Notification.MarkAllRead(ctx, userID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotificationUpdateError, err)
	}
	return &MarkReadOutput{Data: OperationResultData{Message: "notifications marked as read"}}, nil
}
