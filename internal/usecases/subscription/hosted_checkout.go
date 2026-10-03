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
	HostedCheckoutUsecase interface {
		Execute(ctx context.Context, teacherID, bearerToken string, in HostedCheckoutInput) (*HostedCheckoutOutput, apperrors.ApplicationError)
	}

	hostedCheckoutUsecase struct {
		contextFactory appcontext.Factory
	}

	HostedCheckoutInput struct {
		PlanID int `json:"plan_id"`

		PayerEmail string `json:"payer_email"`
	}

	HostedCheckoutData struct {
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
