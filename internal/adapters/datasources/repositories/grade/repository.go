package grade

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Grade) (string, error)
	List(context.Context) ([]domain.Grade, error)
	Get(context.Context, string) (*domain.Grade, error)
	Update(context.Context, string, domain.Grade) error
	Delete(context.Context, string) error
	AddMember(context.Context, domain.GradeMembership) error
	RemoveMember(context.Context, string, string) error
	ListMembers(context.Context, string) ([]domain.UserProfile, error)
	ListUserGrades(context.Context, string) ([]domain.Grade, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, grade domain.Grade) (string, error) {
	query := `
		INSERT INTO grades (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, grade.Name, grade.Description, grade.CreatedBy).Scan(&id)
	return id, err
}

func (r *repository) List(ctx context.Context) ([]domain.Grade, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(description, ''), created_by, created_at
		FROM grades
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grades := []domain.Grade{}
	for rows.Next() {
		var grade domain.Grade
		if err := rows.Scan(&grade.ID, &grade.Name, &grade.Description, &grade.CreatedBy, &grade.CreatedAt); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	return grades, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Grade, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(description, ''), created_by, created_at
		FROM grades
		WHERE id = $1
	`, id)

	var grade domain.Grade
	if err := row.Scan(&grade.ID, &grade.Name, &grade.Description, &grade.CreatedBy, &grade.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}

func (r *repository) Update(ctx context.Context, id string, grade domain.Grade) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE grades SET name = $1, description = $2 WHERE id = $3
	`, grade.Name, grade.Description, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM grades WHERE id = $1`, id)
	return err
}

func (r *repository) AddMember(ctx context.Context, membership domain.GradeMembership) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO grade_memberships (grade_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (grade_id, user_id) DO NOTHING
	`, membership.GradeID, membership.UserID)
	return err
}

func (r *repository) RemoveMember(ctx context.Context, gradeID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM grade_memberships
		WHERE grade_id = $1 AND user_id = $2
	`, gradeID, userID)
	return err
}

func (r *repository) ListMembers(ctx context.Context, gradeID string) ([]domain.UserProfile, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT up.id, up.name, up.email, up.profile_type, up.assistant_base_url, up.assistant_api_key, up.created_at
		FROM user_profiles up
		JOIN grade_memberships gm ON gm.user_id = up.id
		WHERE gm.grade_id = $1
		ORDER BY up.name ASC
	`, gradeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.UserProfile{}
	for rows.Next() {
		var user domain.UserProfile
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.ProfileType,
			&user.AssistantBaseURL,
			&user.AssistantAPIKey,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) ListUserGrades(ctx context.Context, userID string) ([]domain.Grade, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT g.id, g.name, COALESCE(g.description, ''), g.created_by, g.created_at
		FROM grades g
		JOIN grade_memberships gm ON gm.grade_id = g.id
		WHERE gm.user_id = $1
		ORDER BY g.name ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grades := []domain.Grade{}
	for rows.Next() {
		var grade domain.Grade
		if err := rows.Scan(&grade.ID, &grade.Name, &grade.Description, &grade.CreatedBy, &grade.CreatedAt); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	return grades, nil
}
