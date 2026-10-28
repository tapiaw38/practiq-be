package aiconversation

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) AddMessage(ctx context.Context, m domain.AIMessage) (string, error) {
	query := `
		INSERT INTO ai_messages (conversation_id, sender, message_type, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, m.ConversationID, m.Sender, m.MessageType, m.Content).Scan(&id)
	return id, err
}
