package stripe

import (
	"emersonargueta/m/v1/paymentgateway"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
)

var _ paymentgateway.Service = &Service{}

// Service for managing stripe subscriptions.
type Service struct {
	client *Client
}

// NewSubscription for a user with stripe. A new subscription for a
// user can be created so that the user can gain access to
// a domain.
func (s *Service) NewSubscription(email string, planType string, paymentMethodID string) (re *stripe.Subscription, e error) {
	newCustomer, err := s.createCustomer(email)
	if err != nil {
		return nil, ErrStripeNewSubscription
	}

	newSubscription, err := s.createSubscription(planType, paymentMethodID, newCustomer)
	if err != nil {
		return nil, err
	}

	return newSubscription, nil
}
func (s *Service) createCustomer(email string) (*stripe.Customer, error) {
	customerParams := stripe.CustomerParams{Email: &email}
	newCustomer, err := s.client.stripe.Customers.New(&customerParams)
	if err != nil {
		return nil, err
	}

	return newCustomer, nil
}
func (s *Service) createSubscription(planType string, paymentMethodID string, newCustomer *stripe.Customer) (*stripe.Subscription, error) {
	paymentMethod, err := s.attachPaymentMethod(paymentMethodID, newCustomer)
	if err != nil {
		return nil, ErrStripeNewSubscription
	}

	if err := s.updateInvoiceSettingsToDefault(paymentMethod, newCustomer); err != nil {
		return nil, ErrStripeNewSubscription
	}

	subscriptionParams := &stripe.SubscriptionParams{
		Customer: &newCustomer.ID,
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(planType),
			},
		},
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")

	return sub.New(subscriptionParams)
}
func (s *Service) attachPaymentMethod(paymentMethodID string, newCustomer *stripe.Customer) (*stripe.PaymentMethod, error) {

	params := &stripe.PaymentMethodAttachParams{
		Customer: &newCustomer.ID,
	}

	return s.client.stripe.PaymentMethods.Attach(paymentMethodID, params)

}
func (s *Service) updateInvoiceSettingsToDefault(pm *stripe.PaymentMethod, nc *stripe.Customer) error {
	customerParams := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(pm.ID),
		},
	}

	if _, err := s.client.stripe.Customers.Update(nc.ID, customerParams); err != nil {
		return err
	}

	return nil
}
