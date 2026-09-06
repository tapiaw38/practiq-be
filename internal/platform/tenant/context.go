package tenant

import "context"

type contextKey struct{}

func WithSchoolID(ctx context.Context, schoolID string) context.Context {
	return context.WithValue(ctx, contextKey{}, schoolID)
}

func SchoolID(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
