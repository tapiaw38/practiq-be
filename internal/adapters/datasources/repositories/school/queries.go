package school

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) Get(ctx context.Context, id string) (*domain.School, error) {
	var s domain.School
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, kind, billing, COALESCE(created_by, ''), created_at
		FROM schools WHERE id = $1
	`, id).Scan(&s.ID, &s.Name, &s.Kind, &s.Billing, &s.CreatedBy, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) List(ctx context.Context) ([]domain.School, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, kind, billing, COALESCE(created_by, ''), created_at
		FROM schools ORDER BY kind, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schools := []domain.School{}
	for rows.Next() {
		var s domain.School
		if err := rows.Scan(&s.ID, &s.Name, &s.Kind, &s.Billing, &s.CreatedBy, &s.CreatedAt); err != nil {
			return nil, err
		}
		schools = append(schools, s)
	}
	return schools, rows.Err()
}

// Update only touches what a caller named. Kind and billing decide what a
// school allows and what it is charged, so blanking one by omission would
// quietly change both.
func (r *repository) Update(ctx context.Context, id string, s domain.School) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE schools SET
			name    = COALESCE(NULLIF($2, ''), name),
			kind    = COALESCE(NULLIF($3, ''), kind),
			billing = COALESCE(NULLIF($4, ''), billing),
			updated_at = NOW()
		WHERE id = $1
	`, id, s.Name, s.Kind, s.Billing)
	return err
}

func (r *repository) RemoveMember(ctx context.Context, schoolID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM school_members WHERE school_id = $1 AND user_id = $2
	`, schoolID, userID)
	return err
}

func (r *repository) ListMembers(ctx context.Context, schoolID string) ([]domain.SchoolMember, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT school_id, user_id, role, active
		FROM school_members WHERE school_id = $1
		ORDER BY role, created_at
	`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []domain.SchoolMember{}
	for rows.Next() {
		var m domain.SchoolMember
		if err := rows.Scan(&m.SchoolID, &m.UserID, &m.Role, &m.Active); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *repository) CountStudents(ctx context.Context, schoolID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM school_members
		WHERE school_id = $1 AND role = 'student' AND active
	`, schoolID).Scan(&count)
	return count, err
}

// ListStudentsByActivity returns a school's active students, least recently
// active first.
//
// That order is what a downgrade deactivates by. Ordering by seniority instead
// would keep the students who finished months ago and cut the ones sitting in
// class today: a teacher's oldest students are usually their most finished
// ones. Students who never practised sort first, before anyone who did.
func (r *repository) ListStudentsByActivity(ctx context.Context, schoolID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT sm.user_id
		FROM school_members sm
		LEFT JOIN (
			SELECT student_id, MAX(last_practiced_at) AS last_at
			FROM student_topic_progress GROUP BY student_id
		) activity ON activity.student_id = sm.user_id
		WHERE sm.school_id = $1 AND sm.role = 'student' AND sm.active
		ORDER BY activity.last_at ASC NULLS FIRST, sm.created_at ASC
	`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// SetMemberActive is what a downgrade uses to push a student out of a plan and
// what a teacher uses to bring one back. Nothing is deleted: the student keeps
// their account and their history, and stops reaching this school's courses.
func (r *repository) SetMemberActive(ctx context.Context, schoolID, userID string, active bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE school_members SET active = $3
		WHERE school_id = $1 AND user_id = $2
	`, schoolID, userID, active)
	return err
}
