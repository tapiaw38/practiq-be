package notebook

import "context"

func (r *repository) UpdatePageStatement(ctx context.Context, pageID, statementText string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notebook_pages
		SET statement_text = $1
		WHERE id = $2
	`, statementText, pageID)
	return err
}
