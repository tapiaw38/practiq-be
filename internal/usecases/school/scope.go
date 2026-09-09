package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

// Scope is what a caller may see, expressed as the schools they belong to.
//
// Every query that narrows by school goes through here. Writing the rule twice
// is how it drifts: this project has already had a listing that matched on
// enrolments alone and returned an empty home, and a student count that matched
// on assignments alone and undercharged. Both were the same rule written in two
// places, and both times only one copy was fixed.
type Scope struct {
	// Unrestricted is a platform superadmin, who sees every school. It is not
	// "all schools listed" but "do not narrow at all", so a superadmin does not
	// stop seeing things by belonging to none.
	Unrestricted bool
	SchoolIDs    []string
}

// Empty reports a caller who belongs to nothing and is not a superadmin. Their
// queries must return nothing rather than everything — the difference between
// an empty list and a leak.
func (s Scope) Empty() bool {
	return !s.Unrestricted && len(s.SchoolIDs) == 0
}

// ScopeFor resolves which schools a request may read.
func ScopeFor(ctx context.Context, app *appcontext.Context, userID string, isSuperAdmin bool) (Scope, error) {
	if isSuperAdmin {
		return Scope{Unrestricted: true}, nil
	}

	members, err := app.Repositories.School.ListForUser(ctx, userID)
	if err != nil {
		return Scope{}, err
	}

	ids := make([]string, 0, len(members))
	for _, member := range members {
		// A deactivated membership is a student a downgrade pushed out of the
		// plan. They keep their history but stop reaching this school.
		if member.Active {
			ids = append(ids, member.SchoolID)
		}
	}
	return Scope{SchoolIDs: ids}, nil
}
