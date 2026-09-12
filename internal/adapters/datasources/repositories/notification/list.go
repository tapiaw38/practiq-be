package notification

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

const defaultLimit = 50

func (r *repository) ListByUser(ctx context.Context, filter ListFilter) ([]domain.Notification, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 200 {
		limit = defaultLimit
	}

	query := `
		SELECT id, user_id, type, title, COALESCE(body, ''), COALESCE(resource_type, ''),
		       COALESCE(resource_id::text, ''), scheduled_at, read_at, created_at
		FROM notifications
		WHERE user_id = $1
	`
	if filter.UnreadOnly {
		query += ` AND read_at IS NULL`
	}
	query += ` ORDER BY created_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, filter.UserID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []domain.Notification{}
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Body,
			&n.ResourceType,
			&n.ResourceID,
			&n.ScheduledAt,
			&n.ReadAt,
			&n.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (r *repository) CountUnread(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read_at IS NULL
	`, userID).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}
