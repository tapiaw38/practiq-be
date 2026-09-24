// Package pgerr reads the parts of a Postgres error the product needs to act
// on, so a constraint the schema already enforces can be answered with the
// reason instead of a 500.
package pgerr

import (
	"errors"

	"github.com/lib/pq"
)

// uniqueViolation is Postgres' SQLSTATE for a duplicate key.
const uniqueViolation = "23505"

// IsUniqueViolation reports whether err is a duplicate-key rejection from the
// named unique index.
//
// The constraint is named rather than matched loosely: a table usually carries
// more than one unique index, and answering "that name is taken" for whichever
// one happened to fire would be wrong as soon as a second is added.
func IsUniqueViolation(err error, constraint string) bool {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return false
	}
	return string(pqErr.Code) == uniqueViolation && pqErr.Constraint == constraint
}
