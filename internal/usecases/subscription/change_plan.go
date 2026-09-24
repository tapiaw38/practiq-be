package subscription

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	// ChangePlanUsecase moves a paying teacher to another plan.
	//
	// Separate from subscribing because the money works differently: there is
	// already an agreement being charged, and it is restated rather than
	// replaced. Subscribing again would charge a whole new month on top of the
	// one already paid for.
	ChangePlanUsecase interface {
		Execute(ctx context.Context, teacherID, bearerToken string, in ChangePlanInput) (*ChangePlanOutput, apperrors.ApplicationError)
	}

	changePlanUsecase struct {
		contextFactory appcontext.Factory
	}

	ChangePlanInput struct {
		PlanID int `json:"plan_id"`
		// CardTokenID and PaymentMethodID pay the prorated difference, and are
		// only needed when moving up. Single use, like any card token.
		CardTokenID     string `json:"card_token_id"`
		PaymentMethodID string `json:"payment_method_id"`
	}

	ChangePlanData struct {
		// Charged is what was taken now for the rest of the current period.
		// Zero when moving down: the cheaper price starts at renewal.
		Charged float64 `json:"charged"`
	}

	ChangePlanOutput struct {
		Data ChangePlanData `json:"data"`
	}
)

func NewChangePlanUsecase(contextFactory appcontext.Factory) ChangePlanUsecase {
	return &changePlanUsecase{contextFactory: contextFactory}
}

func (u *changePlanUsecase) Execute(ctx context.Context, teacherID, bearerToken string, in ChangePlanInput) (*ChangePlanOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if in.PlanID <= 0 {
		return nil, apperrors.NewBadRequestError("a plan is required")
	}

	// The payer is whoever is asking, and their email comes from auth-api. A
	// card token is not enough to identify anybody: it says how to charge, not
	// who to charge.
	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{teacherID})
	if appErr != nil {
		return nil, appErr
	}
	email := strings.TrimSpace(names[teacherID].Email)
	if email == "" {
		return nil, apperrors.NewBadRequestError("your account has no email to bill")
	}

	change, err := app.Integrations.Payments.ChangePlan(ctx, payments.ChangePlanInput{
		PlanID:          in.PlanID,
		UserID:          teacherID,
		PayerEmail:      email,
		CardTokenID:     in.CardTokenID,
		PaymentMethodID: in.PaymentMethodID,
	})
	if err != nil {
		if rejected, ok := payments.IsRejected(err); ok {
			return nil, apperrors.NewBadRequestError(changePlanMessage(rejected.Code))
		}
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	return &ChangePlanOutput{Data: ChangePlanData{Charged: change.Charged}}, nil
}

func changePlanMessage(code string) string {
	switch {
	case strings.Contains(code, "already_subscribed_to_plan"):
		return "Ya estás en ese plan."
	case strings.Contains(code, "no_live_subscription"):
		return "No tenés una suscripción activa para cambiar."
	case strings.Contains(code, "subscription_not_active"):
		return "Tu suscripción está pausada. Reanudala antes de cambiar de plan."
	case strings.Contains(code, "plan_not_found"):
		return "Ese plan no existe."
	case strings.Contains(code, "rejected"), strings.Contains(code, "token"):
		return "La tarjeta fue rechazada al cobrar la diferencia. Probá con otra."
	default:
		return "No pudimos cambiar el plan. Probá de nuevo en un rato."
	}
}
