package aiconversation

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	CreateConversation(context.Context, domain.AIConversation) (string, error)
	Get(context.Context, string) (*domain.AIConversation, error)
	ListMessages(context.Context, string) ([]domain.AIMessage, error)
	AddMessage(context.Context, domain.AIMessage) (string, error)
	CreateHelpRequest(context.Context, domain.AIHelpRequest) (string, error)
	ListRecentHelpRequests(context.Context, string, string, int) ([]domain.AIHelpRequest, error)
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
