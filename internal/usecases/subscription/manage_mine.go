package subscription

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// Actions a teacher may take on their own subscription.
const (
	ActionPause  = "pause"
	ActionResume = "resume"
	ActionCancel = "cancel"
)

type (
	// ManageMineUsecase pauses, resumes or cancels the asking teacher's
	// subscription.
	//
	// It takes no subscription id. The teacher's subscription is looked up from
	// their own identity, so there is no id for a caller to swap for somebody
	// else's — the one thing that must not be possible here.
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

// currentSubscription finds the one subscription a teacher can act on.
//
// Every subscription is read, not only the live one: a paused subscription is
// not an entitlement, and looking only at entitlements would make resuming
// impossible for exactly the teachers who need it.
func currentSubscription(ctx context.Context, app *appcontext.Context, teacherID string) (*payments.Subscription, apperrors.ApplicationError) {
	subscriptions, err := app.Integrations.Payments.ListSubscriptions(ctx, teacherID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	for _, subscription := range subscriptions {
		// Cancelled agreements are terminal at the gateway and cannot be
		// acted on again.
		if subscription.Status == "cancelled" || subscription.Status == "canceled" {
			continue
		}
		return &subscription, nil
	}
	return nil, apperrors.NewNotFoundError("no subscription to manage")
}
