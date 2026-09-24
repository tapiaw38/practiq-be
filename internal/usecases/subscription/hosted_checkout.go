package subscription

import (
	"context"
	"errors"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/payments"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	// HostedCheckoutUsecase subscribes a teacher who is not giving us a card.
	//
	// Mercado Pago refuses prepaid cards for a monthly charge, and plenty of
	// teachers have no credit card at all. The gateway's own checkout accepts
	// the balance in their Mercado Pago account, which no card form of ours can
	// do, so this hands back the address to send them to.
	HostedCheckoutUsecase interface {
		Execute(ctx context.Context, teacherID, bearerToken string, in HostedCheckoutInput) (*HostedCheckoutOutput, apperrors.ApplicationError)
	}

	hostedCheckoutUsecase struct {
		contextFactory appcontext.Factory
	}

	HostedCheckoutInput struct {
		PlanID int `json:"plan_id"`
		// PayerEmail is the teacher's Mercado Pago address, which is often not
		// the one they signed up to Practiq with. The gateway demands it up
		// front and then refuses the checkout to anyone who logs in with a
		// different account, so guessing it wastes the trip.
		//
		// Unlike the card flow this is safe to take from the request: it names
		// who Mercado Pago asks to authorise the charge, and it cannot charge
		// them — they have to log in and agree first.
		PayerEmail string `json:"payer_email"`
	}

	HostedCheckoutData struct {
		// InitPoint is where the browser has to go. Nothing is charged before
		// the teacher authorises the agreement there.
		InitPoint string `json:"init_point"`
	}

	HostedCheckoutOutput struct {
		Data HostedCheckoutData `json:"data"`
	}
)

func NewHostedCheckoutUsecase(contextFactory appcontext.Factory) HostedCheckoutUsecase {
	return &hostedCheckoutUsecase{contextFactory: contextFactory}
}

func (u *hostedCheckoutUsecase) Execute(ctx context.Context, teacherID, bearerToken string, in HostedCheckoutInput) (*HostedCheckoutOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if in.PlanID <= 0 {
		return nil, apperrors.NewBadRequestError("a plan is required")
	}

	// The account email is only the default. Whoever is asking still comes from
	// the token, so nobody can open a checkout for another teacher's plan.
	email := strings.TrimSpace(in.PayerEmail)
	if email == "" {
		names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{teacherID})
		if appErr != nil {
			return nil, appErr
		}
		email = strings.TrimSpace(names[teacherID].Email)
	}
	if email == "" {
		return nil, apperrors.NewBadRequestError("your account has no email to bill")
	}
	if !strings.Contains(email, "@") || strings.ContainsAny(email, " \t\r\n") {
		return nil, apperrors.NewBadRequestError("ese no parece un email válido")
	}

	hosted, err := app.Integrations.Payments.StartHostedSubscription(ctx, payments.HostedSubscriptionInput{
		PlanID:     in.PlanID,
		UserID:     teacherID,
		PayerEmail: email,
	})
	if err != nil {
		if rejected, ok := payments.IsRejected(err); ok {
			return nil, apperrors.NewBadRequestError(hostedRejectionMessage(rejected.Code))
		}
		return nil, apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	if strings.TrimSpace(hosted.InitPoint) == "" {
		return nil, apperrors.NewApplicationError(
			mappings.SubscriptionUnavailableError,
			errors.New("payments returned no init_point"),
		)
	}
	return &HostedCheckoutOutput{Data: HostedCheckoutData{InitPoint: hosted.InitPoint}}, nil
}

func hostedRejectionMessage(code string) string {
	switch {
	case strings.Contains(code, "already_subscribed"):
		return "Ya tenés una suscripción activa. Para cambiar de plan usá el pago con tarjeta."
	case strings.Contains(code, "plan_not_found"):
		return "Ese plan no existe."
	default:
		return "No pudimos abrir el pago en Mercado Pago. Probá de nuevo en un rato."
	}
}
