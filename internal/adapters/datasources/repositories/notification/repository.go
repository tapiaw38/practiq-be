package notification

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	// Upsert replaces the notification for the same user, type and resource, so
	// rescheduling an event does not stack duplicates.
	Upsert(context.Context, domain.Notification) error
	ListByUser(context.Context, ListFilter) ([]domain.Notification, error)
	CountUnread(context.Context, string) (int, error)
	MarkRead(ctx context.Context, id, userID string) (bool, error)
	MarkAllRead(context.Context, string) error
	DeleteByResource(ctx context.Context, notificationType, resourceID string) error
}

type ListFilter struct {
	UserID     string
	UnreadOnly bool
	Limit      int
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
