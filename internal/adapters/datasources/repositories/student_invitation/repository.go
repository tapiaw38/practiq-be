package studentinvitation

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.StudentInvitation) (*domain.StudentInvitation, error)
	GetByCode(context.Context, string) (*domain.StudentInvitation, error)
	GetActiveByTeacher(context.Context, string) (*domain.StudentInvitation, error)
	Revoke(ctx context.Context, id, teacherID string) error
	// Redeem deja registrado el canje y devuelve si es la primera vez de este
	// alumno con este código.
	Redeem(ctx context.Context, invitationID, studentID string) (bool, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
