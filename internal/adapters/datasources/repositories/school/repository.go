package school

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.School) (string, error)
	CreateWithAdmin(context.Context, domain.School, string) (string, error)

	GetPersonal(ctx context.Context, ownerID string) (*domain.School, error)
	Rename(ctx context.Context, schoolID, name string) error
	AddMember(context.Context, domain.SchoolMember) error

	ListForUser(ctx context.Context, userID string) ([]domain.SchoolMember, error)
	Get(ctx context.Context, id string) (*domain.School, error)
	List(ctx context.Context) ([]domain.School, error)
	Update(ctx context.Context, id string, s domain.School) error
	Close(ctx context.Context, id, closedBy, reason string) error
	Suspend(ctx context.Context, id string) error
	Reopen(ctx context.Context, id string) error
	CountActiveAdmins(ctx context.Context, schoolID string) (int, error)
	RemoveMember(ctx context.Context, schoolID, userID string) error
	ListMembers(ctx context.Context, schoolID string) ([]domain.SchoolMember, error)

	CountStudents(ctx context.Context, schoolID string) (int, error)

	ListStudentsByActivity(ctx context.Context, schoolID string) ([]string, error)
	ListStudentsWithActivity(ctx context.Context, schoolID string) ([]domain.StudentActivity, error)
	SetMemberActive(ctx context.Context, schoolID, userID string, active bool) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
