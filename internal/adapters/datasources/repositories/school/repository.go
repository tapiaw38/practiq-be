package school

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.School) (*domain.School, error)
	List(context.Context, string, bool) ([]domain.School, error)
	Get(context.Context, string) (*domain.School, error)
	ListMembers(context.Context, string) ([]domain.SchoolMembership, error)
	UpsertMember(context.Context, domain.SchoolMembership) error
	DeleteMember(context.Context, string, string) error
	HasAdminAccess(context.Context, string, string) (bool, error)
	IsMember(context.Context, string, string) (bool, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, s domain.School) (*domain.School, error) {
	const query = `
		INSERT INTO schools (name, slug, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, name, slug, is_active, created_by, created_at, updated_at`
	var created domain.School
	err := r.db.QueryRowContext(ctx, query, s.Name, s.Slug, s.CreatedBy).Scan(
		&created.ID, &created.Name, &created.Slug, &created.IsActive,
		&created.CreatedBy, &created.CreatedAt, &created.UpdatedAt,
	)
	return &created, err
}

func (r *repository) List(ctx context.Context, userID string, global bool) ([]domain.School, error) {
	query := `SELECT s.id, s.name, s.slug, s.is_active, s.created_by, s.created_at, s.updated_at
		FROM schools s`
	args := []any{}
	if !global {
		query += ` JOIN school_memberships m ON m.school_id = s.id WHERE m.user_id = $1`
		args = append(args, userID)
	}
	query += ` ORDER BY s.name`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.School, 0)
	for rows.Next() {
		var s domain.School
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.IsActive, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *repository) Get(ctx context.Context, id string) (*domain.School, error) {
	var s domain.School
	err := r.db.QueryRowContext(ctx, `SELECT id, name, slug, is_active, created_by, created_at, updated_at FROM schools WHERE id = $1`, id).
		Scan(&s.ID, &s.Name, &s.Slug, &s.IsActive, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) ListMembers(ctx context.Context, schoolID string) ([]domain.SchoolMembership, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT school_id, user_id, membership_role, profile_type, created_at, updated_at FROM school_memberships WHERE school_id = $1 ORDER BY created_at`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.SchoolMembership, 0)
	for rows.Next() {
		var m domain.SchoolMembership
		if err := rows.Scan(&m.SchoolID, &m.UserID, &m.MembershipRole, &m.ProfileType, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

func (r *repository) UpsertMember(ctx context.Context, m domain.SchoolMembership) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO school_memberships (school_id, user_id, membership_role, profile_type)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (school_id, user_id) DO UPDATE SET
			membership_role = EXCLUDED.membership_role,
			profile_type = EXCLUDED.profile_type,
			updated_at = NOW()`, m.SchoolID, m.UserID, m.MembershipRole, m.ProfileType)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_profiles (id, profile_type, academic_status)
		VALUES ($1, $2, 'active')
		ON CONFLICT (id) DO UPDATE SET profile_type = EXCLUDED.profile_type`, m.UserID, m.ProfileType)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *repository) DeleteMember(ctx context.Context, schoolID, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM school_memberships WHERE school_id = $1 AND user_id = $2`, schoolID, userID)
	return err
}

func (r *repository) HasAdminAccess(ctx context.Context, schoolID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM school_memberships WHERE school_id = $1 AND user_id = $2 AND membership_role = 'admin')`, schoolID, userID).Scan(&exists)
	return exists, err
}

func (r *repository) IsMember(ctx context.Context, schoolID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM school_memberships WHERE school_id = $1 AND user_id = $2)`, schoolID, userID).Scan(&exists)
	return exists, err
}
