package aiconversation

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	CreateConversation(context.Context, domain.AIConversation) (string, error)
	ListMessages(context.Context, string) ([]domain.AIMessage, error)
	AddMessage(context.Context, domain.AIMessage) (string, error)
	CreateHelpRequest(context.Context, domain.AIHelpRequest) (string, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateConversation(ctx context.Context, c domain.AIConversation) (string, error) {
	query := `
		INSERT INTO ai_conversations (student_id, course_id, practice_sheet_id)
		VALUES ($1, NULLIF($2,'')::uuid, NULLIF($3,'')::uuid)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, c.StudentID, c.CourseID, c.PracticeSheetID).Scan(&id)
	return id, err
}

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

func (r *repository) CreateHelpRequest(ctx context.Context, h domain.AIHelpRequest) (string, error) {
	query := `
		INSERT INTO ai_help_requests (student_id, exercise_id, question, ai_response, help_type)
		VALUES ($1, NULLIF($2,'')::uuid, $3, $4, $5)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, h.StudentID, h.ExerciseID, h.Question, h.AIResponse, h.HelpType).Scan(&id)
	return id, err
}
