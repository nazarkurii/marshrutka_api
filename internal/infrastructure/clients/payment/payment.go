package payment

import (
	"maryan_api/config"

	"github.com/stripe/stripe-go"
)

type Payment interface {
	CreateCheckoutSession(amount int64, base, token string) (string, string, error)
}

type paymentImpl struct {
	apiKey string
}

func Init() Payment {
	stripe.Key = config.StripSekretKey()
	return &stripePayment{apiKey: stripe.Key}
}
