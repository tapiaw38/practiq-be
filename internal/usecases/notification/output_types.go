package notification

import "github.com/tapiaw38/practiq-be/internal/domain"

const timeFormat = "2006-01-02T15:04:05Z"

type (
	NotificationData struct {
		ID           string `json:"id"`
		Type         string `json:"type"`
		Title        string `json:"title"`
		Body         string `json:"body,omitempty"`
		ResourceType string `json:"resource_type,omitempty"`
		ResourceID   string `json:"resource_id,omitempty"`
		// ScheduledAt is UTC; the client renders it in the user's timezone.
		ScheduledAt string `json:"scheduled_at,omitempty"`
		Read        bool   `json:"read"`
		CreatedAt   string `json:"created_at"`
	}

	ListData struct {
		Notifications []NotificationData `json:"notifications"`
		UnreadCount   int                `json:"unread_count"`
	}

	OperationResultData struct {
		Message string `json:"message"`
	}
)

func toNotificationData(n domain.Notification) NotificationData {
	data := NotificationData{
		ID:           n.ID,
		Type:         n.Type,
		Title:        n.Title,
		Body:         n.Body,
		ResourceType: n.ResourceType,
		ResourceID:   n.ResourceID,
		Read:         n.ReadAt != nil,
		CreatedAt:    n.CreatedAt.UTC().Format(timeFormat),
	}
	if n.ScheduledAt != nil {
		data.ScheduledAt = n.ScheduledAt.UTC().Format(timeFormat)
	}
	return data
}
