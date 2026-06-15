package teacherstudentassignment

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func (r *repository) listUsers(ctx context.Context, query, id string) ([]domain.UserProfile, error) {
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.UserProfile{}
	for rows.Next() {
		var user domain.UserProfile
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.ProfileType, &user.AcademicStatus, &user.AssistantBaseURL, &user.AssistantAPIKey, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
