package notification

import (
	"context"

	notificationRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/notification"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(context.Context, ListInput) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListInput struct {
		UserID     string
		UnreadOnly bool
		Limit      int
	}

	ListOutput struct {
		Data ListData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, input ListInput) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	notifications, err := app.Repositories.Notification.ListByUser(ctx, notificationRepo.ListFilter{
		UserID:     input.UserID,
		UnreadOnly: input.UnreadOnly,
		Limit:      input.Limit,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotificationListError, err)
	}

	// Always the full unread count, even when the list is filtered or truncated:
	// the badge would be wrong otherwise.
	unread, err := app.Repositories.Notification.CountUnread(ctx, input.UserID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotificationListError, err)
	}

	items := make([]NotificationData, 0, len(notifications))
	for _, n := range notifications {
		items = append(items, toNotificationData(n))
	}
	return &ListOutput{Data: ListData{Notifications: items, UnreadCount: unread}}, nil
}
