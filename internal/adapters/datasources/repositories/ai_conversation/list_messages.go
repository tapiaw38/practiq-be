package aiconversation

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) ListMessages(ctx context.Context, conversationID string) ([]domain.AIMessage, error) {
	query := `
		SELECT id, conversation_id, sender, message_type, COALESCE(content,''), created_at
		FROM ai_messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.AIMessage
	for rows.Next() {
		var m domain.AIMessage
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Sender, &m.MessageType, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
