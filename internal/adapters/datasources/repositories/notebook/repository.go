package notebook

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, n domain.Notebook) (string, error)
	List(ctx context.Context, courseID string) ([]domain.Notebook, error)
	Get(ctx context.Context, id string) (*domain.Notebook, error)

	CreatePage(ctx context.Context, p domain.NotebookPage) (string, error)
	UpdatePage(ctx context.Context, p domain.NotebookPage) error
	DeletePage(ctx context.Context, pageID string) error

	UpsertSubmission(ctx context.Context, s domain.NotebookSubmission) error
	GetSubmission(ctx context.Context, pageID, studentID string) (*domain.NotebookSubmission, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, n domain.Notebook) (string, error) {
	level := n.Level
	if level < 1 {
		level = 1
	}
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO notebooks (course_id, teacher_id, title, description, level)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, n.CourseID, n.TeacherID, n.Title, n.Description, level).Scan(&id)
	return id, err
}

func (r *repository) List(ctx context.Context, courseID string) ([]domain.Notebook, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, course_id, teacher_id, title, description, level, created_at, updated_at
		FROM notebooks WHERE course_id = $1 ORDER BY level ASC, created_at DESC
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notebooks []domain.Notebook
	for rows.Next() {
		var nb domain.Notebook
		if err := rows.Scan(&nb.ID, &nb.CourseID, &nb.TeacherID, &nb.Title, &nb.Description, &nb.Level, &nb.CreatedAt, &nb.UpdatedAt); err != nil {
			return nil, err
		}
		notebooks = append(notebooks, nb)
	}
	return notebooks, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Notebook, error) {
	var nb domain.Notebook
	err := r.db.QueryRowContext(ctx, `
		SELECT id, course_id, teacher_id, title, description, level, created_at, updated_at
		FROM notebooks WHERE id = $1
	`, id).Scan(&nb.ID, &nb.CourseID, &nb.TeacherID, &nb.Title, &nb.Description, &nb.Level, &nb.CreatedAt, &nb.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	pages, err := r.listPages(ctx, id)
	if err != nil {
		return nil, err
	}
	nb.Pages = pages
	return &nb, nil
}

func (r *repository) listPages(ctx context.Context, notebookID string) ([]domain.NotebookPage, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, notebook_id, page_number, title, content_type, content_data, instructions, created_at
		FROM notebook_pages WHERE notebook_id = $1 ORDER BY page_number ASC
	`, notebookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []domain.NotebookPage
	for rows.Next() {
		var p domain.NotebookPage
		if err := rows.Scan(&p.ID, &p.NotebookID, &p.PageNumber, &p.Title, &p.ContentType, &p.ContentData, &p.Instructions, &p.CreatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, nil
}

func (r *repository) CreatePage(ctx context.Context, p domain.NotebookPage) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO notebook_pages (notebook_id, page_number, title, content_type, content_data, instructions)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`, p.NotebookID, p.PageNumber, p.Title, p.ContentType, p.ContentData, p.Instructions).Scan(&id)
	return id, err
}

func (r *repository) UpdatePage(ctx context.Context, p domain.NotebookPage) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_pages
		SET title = $1, content_type = $2, content_data = $3, instructions = $4
		WHERE id = $5
	`, p.Title, p.ContentType, p.ContentData, p.Instructions, p.ID)
	return err
}

func (r *repository) DeletePage(ctx context.Context, pageID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM notebook_pages WHERE id = $1`, pageID)
	return err
}

func (r *repository) UpsertSubmission(ctx context.Context, s domain.NotebookSubmission) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO notebook_submissions (page_id, student_id, canvas_data, answer_text)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (page_id, student_id)
		DO UPDATE SET canvas_data = EXCLUDED.canvas_data, answer_text = EXCLUDED.answer_text, updated_at = NOW()
	`, s.PageID, s.StudentID, s.CanvasData, s.AnswerText)
	return err
}

func (r *repository) GetSubmission(ctx context.Context, pageID, studentID string) (*domain.NotebookSubmission, error) {
	var s domain.NotebookSubmission
	err := r.db.QueryRowContext(ctx, `
		SELECT id, page_id, student_id, canvas_data, answer_text, submitted_at, updated_at
		FROM notebook_submissions WHERE page_id = $1 AND student_id = $2
	`, pageID, studentID).Scan(&s.ID, &s.PageID, &s.StudentID, &s.CanvasData, &s.AnswerText, &s.SubmittedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}
