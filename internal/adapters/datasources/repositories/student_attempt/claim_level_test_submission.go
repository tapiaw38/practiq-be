package studentattempt

import "context"

// ClaimLevelTestSubmission is intentionally an INSERT with a unique key, not
// a read-then-write check. Async submit jobs can run concurrently, and both
// would otherwise observe an empty attempt list before either writes answers.
//
// maxAttempts is how many submissions the sheet allows; the caller passes 1
// for a sheet that sets no limit, which is what a level test allowed before
// limits existed. The cap is part of the same statement for the same reason
// the insert is: two jobs reading the count and then updating it would both
// see room for one more.
//
// Clearing started_at closes the student's window: the next attempt starts its
// own when they open the test again. Sharing one window across attempts made
// the later ones unusable, since a second attempt would begin with whatever
// time the first had left -- often none.
func (r *repository) ClaimLevelTestSubmission(ctx context.Context, studentID, sheetID string, maxAttempts int) (bool, error) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO level_test_submissions (practice_sheet_id, student_id, attempts)
		VALUES ($1::uuid, $2, 1)
		ON CONFLICT (practice_sheet_id, student_id) DO UPDATE
		SET attempts = level_test_submissions.attempts + 1,
		    submitted_at = NOW(),
		    started_at = NULL
		WHERE level_test_submissions.attempts < $3
	`, sheetID, studentID, maxAttempts)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected == 1, nil
}
