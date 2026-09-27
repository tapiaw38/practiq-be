package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

const (
	ActionPause  = "pause"
	ActionResume = "resume"
	ActionCancel = "cancel"
)

type (
	ManageMineUsecase interface {
		Execute(ctx context.Context, teacherID, action string) apperrors.ApplicationError
	}

	manageMineUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewManageMineUsecase(contextFactory appcontext.Factory) ManageMineUsecase {
	return &manageMineUsecase{contextFactory: contextFactory}
}

func (u *manageMineUsecase) Execute(ctx context.Context, teacherID, action string) apperrors.ApplicationError {
	app := u.contextFactory()

	subscription, appErr := currentSubscription(ctx, app, teacherID)
	if appErr != nil {
		return appErr
	}

	var err error
	switch action {
	case ActionPause:
		_, err = app.Integrations.Payments.PauseSubscription(ctx, subscription.ID)
	case ActionResume:
		_, err = app.Integrations.Payments.ResumeSubscription(ctx, subscription.ID)
	case ActionCancel:
		err = app.Integrations.Payments.CancelSubscription(ctx, subscription.ID)
	default:
		return apperrors.NewBadRequestError("unknown subscription action")
	}
	if err != nil {
		return apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	return nil
}

func currentSubscription(ctx context.Context, app *appcontext.Context, teacherID string) (*payments.Subscription, apperrors.ApplicationError) {
	subscriptions, err := app.Integrations.Payments.ListSubscriptions(ctx, teacherID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	for _, subscription := range subscriptions {

		if subscription.Status == "cancelled" || subscription.Status == "canceled" {
			continue
		}
		return &subscription, nil
	}
	return nil, apperrors.NewNotFoundError("no subscription to manage")
}
