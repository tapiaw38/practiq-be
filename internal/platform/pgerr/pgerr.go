package pgerr

import (
	"errors"

	"github.com/lib/pq"
)

const uniqueViolation = "23505"

func IsUniqueViolation(err error, constraint string) bool {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return false
	}
	return string(pqErr.Code) == uniqueViolation && pqErr.Constraint == constraint
}
