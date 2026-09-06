package studentattempt

import (
	"context"
	"github.com/tapiaw38/practiq-be/internal/platform/tenant"
)

func (r *repository) DeleteBySheet(ctx context.Context, studentID, sheetID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM student_attempts
		WHERE student_id = $1 AND practice_sheet_id = $2::uuid AND ($3 = '' OR school_id = NULLIF($3, '')::uuid)
	`, studentID, sheetID, tenant.SchoolID(ctx))
	return err
}
