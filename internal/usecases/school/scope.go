package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type Scope struct {
	Unrestricted bool
	SchoolIDs    []string
}

func (s Scope) Empty() bool {
	return !s.Unrestricted && len(s.SchoolIDs) == 0
}

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

		if member.Active {
			ids = append(ids, member.SchoolID)
		}
	}
	return Scope{SchoolIDs: ids}, nil
}

func ScopeForSchool(ctx context.Context, app *appcontext.Context, userID string, isSuperAdmin bool, schoolID string) (Scope, error) {
	if schoolID == "" {
		return ScopeFor(ctx, app, userID, isSuperAdmin)
	}
	if isSuperAdmin {
		return Scope{SchoolIDs: []string{schoolID}}, nil
	}
	members, err := app.Repositories.School.ListForUser(ctx, userID)
	if err != nil {
		return Scope{}, err
	}
	for _, member := range members {
		if member.SchoolID == schoolID && member.Active {
			return Scope{SchoolIDs: []string{schoolID}}, nil
		}
	}
	return Scope{SchoolIDs: []string{}}, nil
}
