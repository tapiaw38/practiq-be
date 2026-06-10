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
	Update(ctx context.Context, id string, n domain.Notebook) error
	Delete(ctx context.Context, id string) error

	CreatePage(ctx context.Context, p domain.NotebookPage) (string, error)
	GetPage(ctx context.Context, pageID string) (*domain.NotebookPage, error)
	UpdatePage(ctx context.Context, p domain.NotebookPage) error
	DeletePage(ctx context.Context, pageID string) error

	UpsertSubmission(ctx context.Context, s domain.NotebookSubmission) error
	GetSubmission(ctx context.Context, pageID, studentID string) (*domain.NotebookSubmission, error)
	GetSubmissionByID(ctx context.Context, id string) (*domain.NotebookSubmission, error)
	GetFullSubmissionByID(ctx context.Context, id string) (*domain.NotebookSubmissionFull, error)
	ListSubmissions(ctx context.Context, filter SubmissionFilter) ([]domain.NotebookSubmissionFull, error)
	UpdateSubmissionAIReview(ctx context.Context, id string, recognizedText string, isCorrect *bool, feedback string) error
	UpdateSubmissionTeacherReview(ctx context.Context, id string, isCorrect bool, feedback string) error
}

type SubmissionFilter struct {
	NotebookID string
	StudentID  string
	CourseID   string
	Reviewed   string
	TeacherID  string
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(notebooks) == 0 {
		return notebooks, nil
	}

	ids := make([]string, len(notebooks))
	for i, nb := range notebooks {
		ids[i] = nb.ID
	}

	pRows, err := r.db.QueryContext(ctx, `
		SELECT id, notebook_id, page_number, title, content_type, content_data, instructions, created_at
		FROM notebook_pages WHERE notebook_id = ANY($1::uuid[]) ORDER BY notebook_id, page_number ASC
	`, "{"+joinIDs(ids)+"}")
	if err != nil {
		return nil, err
	}
	defer pRows.Close()

	index := make(map[string]int, len(notebooks))
	for i, nb := range notebooks {
		index[nb.ID] = i
	}
	for pRows.Next() {
		var p domain.NotebookPage
		if err := pRows.Scan(&p.ID, &p.NotebookID, &p.PageNumber, &p.Title, &p.ContentType, &p.ContentData, &p.Instructions, &p.CreatedAt); err != nil {
			return nil, err
		}
		if i, ok := index[p.NotebookID]; ok {
			notebooks[i].Pages = append(notebooks[i].Pages, p)
		}
	}
	return notebooks, nil
}

func joinIDs(ids []string) string {
	result := ""
	for i, id := range ids {
		if i > 0 {
			result += ","
		}
		result += id
	}
	return result
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

func (r *repository) Update(ctx context.Context, id string, n domain.Notebook) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebooks SET title = $1, description = $2, updated_at = NOW() WHERE id = $3
	`, n.Title, n.Description, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM notebooks WHERE id = $1`, id)
	return err
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

func (r *repository) GetPage(ctx context.Context, pageID string) (*domain.NotebookPage, error) {
	var p domain.NotebookPage
	err := r.db.QueryRowContext(ctx, `
		SELECT id, notebook_id, page_number, title, content_type, content_data, instructions, created_at
		FROM notebook_pages WHERE id = $1
	`, pageID).Scan(&p.ID, &p.NotebookID, &p.PageNumber, &p.Title, &p.ContentType, &p.ContentData, &p.Instructions, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
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
		INSERT INTO notebook_submissions (page_id, student_id, canvas_data, answer_text, ai_recognized_text, ai_is_correct, ai_feedback, ai_reviewed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (page_id, student_id)
		DO UPDATE SET canvas_data = EXCLUDED.canvas_data, answer_text = EXCLUDED.answer_text, ai_recognized_text = EXCLUDED.ai_recognized_text, ai_is_correct = EXCLUDED.ai_is_correct, ai_feedback = EXCLUDED.ai_feedback, ai_reviewed_at = EXCLUDED.ai_reviewed_at, updated_at = NOW()
	`, s.PageID, s.StudentID, s.CanvasData, s.AnswerText, s.AIRecognizedText, s.AIIsCorrect, s.AIFeedback, s.AIReviewedAt)
	return err
}

func (r *repository) GetSubmission(ctx context.Context, pageID, studentID string) (*domain.NotebookSubmission, error) {
	var s domain.NotebookSubmission
	err := r.db.QueryRowContext(ctx, `
		SELECT id, page_id, student_id, canvas_data, answer_text, COALESCE(ai_recognized_text,''), ai_is_correct, ai_feedback, ai_reviewed_at,
		       teacher_is_correct, COALESCE(teacher_feedback,''), teacher_reviewed_at, submitted_at, updated_at
		FROM notebook_submissions WHERE page_id = $1 AND student_id = $2
	`, pageID, studentID).Scan(&s.ID, &s.PageID, &s.StudentID, &s.CanvasData, &s.AnswerText, &s.AIRecognizedText, &s.AIIsCorrect, &s.AIFeedback, &s.AIReviewedAt, &s.TeacherIsCorrect, &s.TeacherFeedback, &s.TeacherReviewedAt, &s.SubmittedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (r *repository) GetSubmissionByID(ctx context.Context, id string) (*domain.NotebookSubmission, error) {
	var s domain.NotebookSubmission
	err := r.db.QueryRowContext(ctx, `
		SELECT id, page_id, student_id, canvas_data, answer_text, COALESCE(ai_recognized_text,''), ai_is_correct, ai_feedback, ai_reviewed_at,
		       teacher_is_correct, COALESCE(teacher_feedback,''), teacher_reviewed_at, submitted_at, updated_at
		FROM notebook_submissions WHERE id = $1
	`, id).Scan(&s.ID, &s.PageID, &s.StudentID, &s.CanvasData, &s.AnswerText, &s.AIRecognizedText, &s.AIIsCorrect, &s.AIFeedback, &s.AIReviewedAt, &s.TeacherIsCorrect, &s.TeacherFeedback, &s.TeacherReviewedAt, &s.SubmittedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) GetFullSubmissionByID(ctx context.Context, id string) (*domain.NotebookSubmissionFull, error) {
	var s domain.NotebookSubmissionFull
	err := r.db.QueryRowContext(ctx, `
		SELECT ns.id, ns.page_id, ns.student_id, ns.canvas_data, ns.answer_text,
		       COALESCE(ns.ai_recognized_text,''), ns.ai_is_correct, COALESCE(ns.ai_feedback,''), ns.ai_reviewed_at,
		       ns.teacher_is_correct, COALESCE(ns.teacher_feedback,''), ns.teacher_reviewed_at,
		       ns.submitted_at, ns.updated_at,
		       COALESCE(up.name,''), COALESCE(up.email,''),
		       n.id, n.title, COALESCE(np.title,''), np.page_number, n.course_id::text, n.teacher_id
		FROM notebook_submissions ns
		JOIN notebook_pages np ON np.id = ns.page_id
		JOIN notebooks n ON n.id = np.notebook_id
		LEFT JOIN user_profiles up ON up.id = ns.student_id
		WHERE ns.id = $1
	`, id).Scan(
		&s.ID, &s.PageID, &s.StudentID, &s.CanvasData, &s.AnswerText,
		&s.AIRecognizedText, &s.AIIsCorrect, &s.AIFeedback, &s.AIReviewedAt,
		&s.TeacherIsCorrect, &s.TeacherFeedback, &s.TeacherReviewedAt,
		&s.SubmittedAt, &s.UpdatedAt,
		&s.StudentName, &s.StudentEmail,
		&s.NotebookID, &s.NotebookTitle, &s.PageTitle, &s.PageNumber, &s.CourseID, &s.TeacherID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) ListSubmissions(ctx context.Context, filter SubmissionFilter) ([]domain.NotebookSubmissionFull, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ns.id, ns.page_id, ns.student_id, ns.canvas_data, ns.answer_text,
		       COALESCE(ns.ai_recognized_text,''), ns.ai_is_correct, COALESCE(ns.ai_feedback,''), ns.ai_reviewed_at,
		       ns.teacher_is_correct, COALESCE(ns.teacher_feedback,''), ns.teacher_reviewed_at,
		       ns.submitted_at, ns.updated_at,
		       COALESCE(up.name,''), COALESCE(up.email,''),
		       n.id, n.title, COALESCE(np.title,''), np.page_number, n.course_id::text, n.teacher_id
		FROM notebook_submissions ns
		JOIN notebook_pages np ON np.id = ns.page_id
		JOIN notebooks n ON n.id = np.notebook_id
		LEFT JOIN user_profiles up ON up.id = ns.student_id
		WHERE ($1 = '' OR n.id::text = $1)
		  AND ($2 = '' OR ns.student_id = $2)
		  AND ($3 = '' OR n.course_id::text = $3)
		  AND ($4 = '' OR ($4 = 'reviewed' AND ns.ai_reviewed_at IS NOT NULL) OR ($4 = 'unreviewed' AND ns.ai_reviewed_at IS NULL))
		  AND ($5 = '' OR n.teacher_id = $5)
		ORDER BY ns.submitted_at DESC
	`, filter.NotebookID, filter.StudentID, filter.CourseID, filter.Reviewed, filter.TeacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.NotebookSubmissionFull
	for rows.Next() {
		var s domain.NotebookSubmissionFull
		if err := rows.Scan(
			&s.ID, &s.PageID, &s.StudentID, &s.CanvasData, &s.AnswerText,
			&s.AIRecognizedText, &s.AIIsCorrect, &s.AIFeedback, &s.AIReviewedAt,
			&s.TeacherIsCorrect, &s.TeacherFeedback, &s.TeacherReviewedAt,
			&s.SubmittedAt, &s.UpdatedAt,
			&s.StudentName, &s.StudentEmail,
			&s.NotebookID, &s.NotebookTitle, &s.PageTitle, &s.PageNumber, &s.CourseID, &s.TeacherID,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *repository) UpdateSubmissionAIReview(ctx context.Context, id string, recognizedText string, isCorrect *bool, feedback string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_submissions
		SET ai_recognized_text = $1, ai_is_correct = $2, ai_feedback = $3, ai_reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $4
	`, recognizedText, isCorrect, feedback, id)
	return err
}

func (r *repository) UpdateSubmissionTeacherReview(ctx context.Context, id string, isCorrect bool, feedback string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_submissions
		SET teacher_is_correct = $1, teacher_feedback = $2, teacher_reviewed_at = NOW(), updated_at = NOW()
		WHERE id = $3
	`, isCorrect, feedback, id)
	return err
}
