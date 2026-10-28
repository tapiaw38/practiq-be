package notification

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Upsert(context.Context, domain.Notification) error
	ListByUser(context.Context, ListFilter) ([]domain.Notification, error)
	CountUnread(context.Context, string) (int, error)
	MarkRead(ctx context.Context, id, userID string) (bool, error)
	MarkAllRead(context.Context, string) error

	Delete(ctx context.Context, id, userID string) (bool, error)
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
