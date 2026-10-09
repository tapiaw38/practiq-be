package notification

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	DeleteUsecase interface {
		Execute(ctx context.Context, id, userID string) (*MarkReadOutput, apperrors.ApplicationError)
	}

	deleteUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{contextFactory: contextFactory}
}

func (u *deleteUsecase) Execute(ctx context.Context, id, userID string) (*MarkReadOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if _, err := app.Repositories.Notification.Delete(ctx, id, userID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotificationUpdateError, err)
	}
	return &MarkReadOutput{Data: OperationResultData{Message: "notification deleted"}}, nil
}
