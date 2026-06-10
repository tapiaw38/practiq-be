package learningstrategy

import (
	"context"
	"database/sql"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type Repository interface {
	List(context.Context) ([]domain.LearningStrategy, error)
	Get(context.Context, string) (*domain.LearningStrategy, error)
	GetByCode(context.Context, string) (*domain.LearningStrategy, error)
	Create(context.Context, domain.LearningStrategy) (string, error)
	Update(context.Context, string, domain.LearningStrategy) error
	Delete(context.Context, string) error

	// Course learning strategies
	ListByCourse(ctx context.Context, courseID string) ([]domain.CourseLearningStrategy, error)
	AssignToCourse(ctx context.Context, cls domain.CourseLearningStrategy) (string, error)
	GetCourseStrategy(ctx context.Context, id string) (*domain.CourseLearningStrategy, error)
	UnassignFromCourse(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context) ([]domain.LearningStrategy, error) {
	query := `SELECT id, name, code, COALESCE(description,''), status, created_at FROM learning_strategies WHERE status='active' ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var strategies []domain.LearningStrategy
	for rows.Next() {
		var s domain.LearningStrategy
		if err := rows.Scan(&s.ID, &s.Name, &s.Code, &s.Description, &s.Status, &s.CreatedAt); err != nil {
			return nil, err
		}
		strategies = append(strategies, s)
	}
	return strategies, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.LearningStrategy, error) {
	query := `SELECT id, name, code, COALESCE(description,''), status, created_at FROM learning_strategies WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var s domain.LearningStrategy
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.Description, &s.Status, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) GetByCode(ctx context.Context, code string) (*domain.LearningStrategy, error) {
	query := `SELECT id, name, code, COALESCE(description,''), status, created_at FROM learning_strategies WHERE code = $1`
	row := r.db.QueryRowContext(ctx, query, code)
	var s domain.LearningStrategy
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.Description, &s.Status, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (r *repository) Create(ctx context.Context, s domain.LearningStrategy) (string, error) {
	query := `
		INSERT INTO learning_strategies (name, code, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	status := s.Status
	if status == "" {
		status = "active"
	}
	var id string
	err := r.db.QueryRowContext(ctx, query, s.Name, s.Code, s.Description, status).Scan(&id)
	return id, err
}

func (r *repository) Update(ctx context.Context, id string, s domain.LearningStrategy) error {
	query := `
		UPDATE learning_strategies
		SET name = $1, code = $2, description = $3, status = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query, s.Name, s.Code, s.Description, s.Status, id)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM learning_strategies WHERE id = $1`, id)
	return err
}

func (r *repository) ListByCourse(ctx context.Context, courseID string) ([]domain.CourseLearningStrategy, error) {
	query := `
		SELECT cls.id, cls.course_id, cls.strategy_id, cls.is_default, COALESCE(cls.config::text, '{}'), cls.created_at,
		       ls.name, ls.code, COALESCE(ls.description,'')
		FROM course_learning_strategies cls
		JOIN learning_strategies ls ON ls.id = cls.strategy_id
		WHERE cls.course_id = $1
		ORDER BY cls.created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.CourseLearningStrategy
	for rows.Next() {
		var cls domain.CourseLearningStrategy
		var configJSON string
		if err := rows.Scan(&cls.ID, &cls.CourseID, &cls.StrategyID, &cls.IsDefault, &configJSON, &cls.CreatedAt,
			&cls.StrategyName, &cls.StrategyCode, &cls.StrategyDescription); err != nil {
			return nil, err
		}
		cls.Config = configJSON
		list = append(list, cls)
	}
	return list, nil
}

func (r *repository) AssignToCourse(ctx context.Context, cls domain.CourseLearningStrategy) (string, error) {
	query := `
		INSERT INTO course_learning_strategies (course_id, strategy_id, is_default, config)
		VALUES ($1, $2, $3, $4::jsonb)
		ON CONFLICT (course_id, strategy_id) DO UPDATE SET is_default = EXCLUDED.is_default, config = EXCLUDED.config
		RETURNING id
	`
	config := cls.Config
	if config == "" {
		config = "{}"
	}
	var id string
	err := r.db.QueryRowContext(ctx, query, cls.CourseID, cls.StrategyID, cls.IsDefault, config).Scan(&id)
	return id, err
}

func (r *repository) GetCourseStrategy(ctx context.Context, id string) (*domain.CourseLearningStrategy, error) {
	query := `
		SELECT cls.id, cls.course_id, cls.strategy_id, cls.is_default, COALESCE(cls.config::text, '{}'), cls.created_at,
		       ls.name, ls.code, COALESCE(ls.description,'')
		FROM course_learning_strategies cls
		JOIN learning_strategies ls ON ls.id = cls.strategy_id
		WHERE cls.id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var cls domain.CourseLearningStrategy
	var configJSON string
	err := row.Scan(&cls.ID, &cls.CourseID, &cls.StrategyID, &cls.IsDefault, &configJSON, &cls.CreatedAt,
		&cls.StrategyName, &cls.StrategyCode, &cls.StrategyDescription)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	cls.Config = configJSON
	return &cls, nil
}

func (r *repository) UnassignFromCourse(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM course_learning_strategies WHERE id = $1`, id)
	return err
}
