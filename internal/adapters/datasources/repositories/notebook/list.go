package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) List(ctx context.Context, courseID string) ([]domain.Notebook, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT n.id, n.course_id, n.teacher_id, n.title, n.description, n.level, n.created_at, n.updated_at
		FROM notebooks n
		JOIN courses c ON c.id = n.course_id
		WHERE n.course_id = $1 AND n.deleted_at IS NULL AND c.deleted_at IS NULL AND ($2 = '' OR c.school_id = NULLIF($2, '')::uuid)
		ORDER BY n.level ASC, n.created_at DESC
	`, courseID, tenant.SchoolID(ctx))
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

	// content_data is deliberately absent: it is the page canvas as a base64
	// data URL, hundreds of KB each, and a listing only ever shows titles and
	// counts. Shipping it made this the slowest endpoint on the dashboard by an
	// order of magnitude. Get(id) still returns it for the screens that draw a
	// page.
	pRows, err := r.db.QueryContext(ctx, `
		SELECT np.id, np.notebook_id, np.page_number, np.title, np.content_type, np.instructions, np.created_at
		FROM notebook_pages np JOIN notebooks n ON n.id = np.notebook_id JOIN courses c ON c.id = n.course_id
		WHERE np.notebook_id = ANY($1::uuid[]) AND ($2 = '' OR c.school_id = NULLIF($2, '')::uuid)
		ORDER BY np.notebook_id, np.page_number ASC
	`, "{"+joinIDs(ids)+"}", tenant.SchoolID(ctx))
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
		if err := pRows.Scan(&p.ID, &p.NotebookID, &p.PageNumber, &p.Title, &p.ContentType, &p.Instructions, &p.CreatedAt); err != nil {
			return nil, err
		}
		if i, ok := index[p.NotebookID]; ok {
			notebooks[i].Pages = append(notebooks[i].Pages, p)
		}
	}
	return notebooks, nil
}
