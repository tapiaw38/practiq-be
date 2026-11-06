package grade

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Grade) (string, error)

	List(ctx context.Context, schoolIDs []string) ([]domain.Grade, error)
	Get(context.Context, string) (*domain.Grade, error)
	Update(context.Context, string, domain.Grade) error
	Delete(context.Context, string) error
	AddMember(context.Context, domain.GradeMembership) error
	RemoveMember(context.Context, string, string) error
	ListMembers(context.Context, string) ([]domain.UserProfile, error)
	ListUserGrades(context.Context, string) ([]domain.Grade, error)

	ListGradesByUsers(context.Context, []string) (map[string][]domain.Grade, error)
}
type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
