package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

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
