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
	// SubscribeUsecase puts the asking teacher on a plan.
	//
	// The card never reaches this service. The browser turns it into a
	// single-use token against the gateway, and that token is all this
	// forwards — which is what keeps card data out of our logs, our database
	// and our compliance scope.
	SubscribeUsecase interface {
		Execute(ctx context.Context, teacherID, bearerToken string, in SubscribeInput) apperrors.ApplicationError
	}

	subscribeUsecase struct {
		contextFactory appcontext.Factory
	}

	SubscribeInput struct {
		PlanID int `json:"plan_id"`
		// CardTokenID is single use and expires quickly. It is not a card
		// number and cannot be charged again on its own.
		CardTokenID string `json:"card_token_id"`
	}

	// CheckoutConfigData carries what the browser needs to tokenise a card.
	CheckoutConfigData struct {
		// PublicKey is public by design. The secret half stays in the payments
		// service and is never served.
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

	// The payer is whoever is asking, and their email comes from auth-api
	// rather than the request. Accepting an email from the body would let a
	// teacher bill a subscription to somebody else's address; accepting an id
	// would let them subscribe another account entirely.
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, []string{teacherID})
	if err != nil {
		return apperrors.NewApplicationError(mappings.ProfileGetError, err)
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
		return apperrors.NewApplicationError(mappings.SubscriptionUnavailableError, err)
	}
	return nil
}

// CheckoutConfig serves the gateway's public key to the browser.
func CheckoutConfig(publicKey string) *CheckoutConfigOutput {
	return &CheckoutConfigOutput{Data: CheckoutConfigData{PublicKey: publicKey}}
}
