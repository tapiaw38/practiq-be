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
	SubscribeUsecase interface {
		Execute(ctx context.Context, teacherID, bearerToken string, in SubscribeInput) apperrors.ApplicationError
	}

	subscribeUsecase struct {
		contextFactory appcontext.Factory
	}

	SubscribeInput struct {
		PlanID int `json:"plan_id"`

		CardTokenID string `json:"card_token_id"`
	}

	CheckoutConfigData struct {
		PublicKey string `json:"public_key"`
	}

	CheckoutConfigOutput struct {
		Data CheckoutConfigData `json:"data"`
	}
)

func NewSubscribeUsecase(contextFactory appcontext.Factory) SubscribeUsecase {
	return &subscribeUsecase{contextFactory: contextFactory}
}

func (u *subscribeUsecase) Execute(ctx context.Context, teacherID, bearerToken string, in SubscribeInput) apperrors.ApplicationError {
	app := u.contextFactory()

	if in.PlanID <= 0 || strings.TrimSpace(in.CardTokenID) == "" {
		return apperrors.NewBadRequestError("a plan and a card token are required")
	}

	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{teacherID})
	if appErr != nil {
		return appErr
	}
	email := strings.TrimSpace(names[teacherID].Email)
	if email == "" {
		return apperrors.NewBadRequestError("your account has no email to bill")
	}

	if _, err := app.Integrations.Payments.CreateSubscription(ctx, payments.SubscriptionInput{
		PlanID:      in.PlanID,
		UserID:      teacherID,
		PayerEmail:  email,
		CardTokenID: in.CardTokenID,
	}); err != nil {
		if rejected, ok := payments.IsRejected(err); ok {
			return apperrors.NewBadRequestError(subscriptionRejectionMessage(rejected.Code))
		}
		return apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	return nil
}

func subscriptionRejectionMessage(code string) string {
	switch {
	case strings.Contains(code, "invalid_user"):
		return "Usá una cuenta de Mercado Pago distinta a la cuenta que recibe el cobro."
	case strings.Contains(code, "token"):
		return "No pudimos validar la tarjeta. Volvé a cargar los datos e intentá de nuevo."
	case strings.Contains(code, "rejected"):

		return "La tarjeta fue rechazada. Las prepagas no sirven para pagos mensuales: usá una de crédito o débito."
	default:
		return "Mercado Pago no pudo autorizar la suscripción. Revisá los datos o probá otra tarjeta."
	}
}

func CheckoutConfig(publicKey string) *CheckoutConfigOutput {
	return &CheckoutConfigOutput{Data: CheckoutConfigData{PublicKey: publicKey}}
}
