package teacherstudentassignment

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Assign(context.Context, domain.TeacherStudentAssignment) error
	Unassign(context.Context, string, string) error
	ListTeachers(context.Context, string) ([]domain.UserProfile, error)
	ListStudents(context.Context, string) ([]domain.UserProfile, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Assign(ctx context.Context, assignment domain.TeacherStudentAssignment) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO teacher_student_assignments (teacher_id, student_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (teacher_id, student_id) DO UPDATE SET status = EXCLUDED.status
	`, assignment.TeacherID, assignment.StudentID, assignment.Status)
	return err
}

func (r *repository) Unassign(ctx context.Context, teacherID, studentID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM teacher_student_assignments
		WHERE teacher_id = $1 AND student_id = $2
	`, teacherID, studentID)
	return err
}

func (r *repository) ListTeachers(ctx context.Context, studentID string) ([]domain.UserProfile, error) {
	return r.listUsers(ctx, `
		SELECT up.id, up.name, up.email, up.profile_type, up.academic_status, up.assistant_base_url, up.assistant_api_key, up.created_at
		FROM user_profiles up
		JOIN teacher_student_assignments tsa ON tsa.teacher_id = up.id
		WHERE tsa.student_id = $1 AND tsa.status = 'active'
		ORDER BY up.name ASC
	`, studentID)
}

func (r *repository) ListStudents(ctx context.Context, teacherID string) ([]domain.UserProfile, error) {
	return r.listUsers(ctx, `
		SELECT up.id, up.name, up.email, up.profile_type, up.academic_status, up.assistant_base_url, up.assistant_api_key, up.created_at
		FROM user_profiles up
		JOIN teacher_student_assignments tsa ON tsa.student_id = up.id
		WHERE tsa.teacher_id = $1 AND tsa.status = 'active'
		ORDER BY up.name ASC
	`, teacherID)
}

func (r *repository) listUsers(ctx context.Context, query, id string) ([]domain.UserProfile, error) {
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.UserProfile{}
	for rows.Next() {
		var user domain.UserProfile
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.ProfileType, &user.AcademicStatus, &user.AssistantBaseURL, &user.AssistantAPIKey, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
