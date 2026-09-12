package school

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.School) (string, error)
	// GetPersonal returns the school a teacher owns, or nil. Personal schools
	// are the only ones with a single owner, so this is the only lookup that
	// can be answered by user alone.
	GetPersonal(ctx context.Context, ownerID string) (*domain.School, error)
	Rename(ctx context.Context, schoolID, name string) error
	AddMember(context.Context, domain.SchoolMember) error
	// ListForUser returns every school a user belongs to with the role they
	// hold in each. Someone can administer one and teach at another.
	ListForUser(ctx context.Context, userID string) ([]domain.SchoolMember, error)
	Get(ctx context.Context, id string) (*domain.School, error)
	List(ctx context.Context) ([]domain.School, error)
	Update(ctx context.Context, id string, s domain.School) error
	Close(ctx context.Context, id, closedBy, reason string) error
	Reopen(ctx context.Context, id string) error
	CountActiveAdmins(ctx context.Context, schoolID string) (int, error)
	RemoveMember(ctx context.Context, schoolID, userID string) error
	ListMembers(ctx context.Context, schoolID string) ([]domain.SchoolMember, error)
	// CountStudents counts the active students of one school. The limit is a
	// school's, not a teacher's: a teacher with their own school who also
	// teaches at an institution must not have those students charged to their
	// personal plan.
	CountStudents(ctx context.Context, schoolID string) (int, error)
	// ListStudentsByActivity orders a school's active students least recently
	// active first, which is the order a downgrade deactivates by.
	ListStudentsByActivity(ctx context.Context, schoolID string) ([]string, error)
	SetMemberActive(ctx context.Context, schoolID, userID string, active bool) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
