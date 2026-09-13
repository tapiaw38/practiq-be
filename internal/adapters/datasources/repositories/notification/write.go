package notification

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Upsert(ctx context.Context, n domain.Notification) error {
	// Rescheduling refreshes the contents and clears read_at, so a student who
	// already dismissed the old date sees the new one.
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO notifications (user_id, type, title, body, resource_type, resource_id, scheduled_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, type, resource_id) WHERE resource_id IS NOT NULL
		DO UPDATE SET
			title = EXCLUDED.title,
			body = EXCLUDED.body,
			scheduled_at = EXCLUDED.scheduled_at,
			read_at = NULL,
			created_at = NOW()
	`, n.UserID, n.Type, n.Title, n.Body, nullable(n.ResourceType), nullable(n.ResourceID), n.ScheduledAt)
	return err
}

func (r *repository) MarkRead(ctx context.Context, id, userID string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE notifications SET read_at = NOW()
		WHERE id = $1 AND user_id = $2 AND read_at IS NULL
	`, id, userID)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	return affected > 0, err
}

func (r *repository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notifications SET read_at = NOW()
		WHERE user_id = $1 AND read_at IS NULL
	`, userID)
	return err
}

func (r *repository) Delete(ctx context.Context, id, userID string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM notifications WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	return affected > 0, err
}

// DeleteByResource drops notifications whose event no longer exists or was
// unscheduled.
func (r *repository) DeleteByResource(ctx context.Context, notificationType, resourceID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM notifications WHERE type = $1 AND resource_id = $2
	`, notificationType, resourceID)
	return err
}

// nullable keeps empty strings out of columns the unique index treats as NULL.
func nullable(value string) any {
	if value == "" {
		return nil
	}
	return value
}
