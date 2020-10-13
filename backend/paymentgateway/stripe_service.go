package paymentgateway

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"github.com/stripe/stripe-go/sub"
)

var _ Processes = &stripeservice{}

type stripeservice struct {
	client *Client
	*client.API
}

// Initialize the underlying stripe client.
func (s *stripeservice) Initialize() {
	apiKey := s.client.config.PaymentGateway.APIKey
	stripe := client.New(apiKey, nil)

	s.API = stripe
}

// NewSubscription for a user with stripe. A new subscription for a
// user can be created so that the user can gain access to
// a domain.
func (s *stripeservice) NewSubscription(details map[string]string) (e error) {
	if e = s.validateDetails(details); e != nil {
		return e
	}

	email := details["email"]
	planType := details["plantype"]
	paymentMethodID := details["paymentmethodid"]

	newCustomer, e := s.createCustomer(email)
	if e != nil {
		return nil
	}

	_, e = s.createSubscription(planType, paymentMethodID, newCustomer)

	return e
}

// NewPayment made by customer.
func (s *stripeservice) NewPayment(details map[string]string) (e error) {
	return e
}

func (s *stripeservice) validateDetails(details map[string]string) (e error) {
	if _, ok := details["email"]; !ok {
		return ErrStripeNewSubscription
	}
	if _, ok := details["plantype"]; !ok {
		return ErrStripeNewSubscription
	}
	if _, ok := details["paymentmethodid"]; !ok {
		return ErrStripeNewSubscription
	}
	return e
}
func (s *stripeservice) createCustomer(email string) (*stripe.Customer, error) {
	customerParams := stripe.CustomerParams{Email: &email}
	newCustomer, err := s.Customers.New(&customerParams)
	if err != nil {
		return nil, err
	}

	return newCustomer, nil
}
func (s *stripeservice) createSubscription(planType string, paymentMethodID string, newCustomer *stripe.Customer) (*stripe.Subscription, error) {
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
func (s *stripeservice) attachPaymentMethod(paymentMethodID string, newCustomer *stripe.Customer) (*stripe.PaymentMethod, error) {

	params := &stripe.PaymentMethodAttachParams{
		Customer: &newCustomer.ID,
	}

	return s.PaymentMethods.Attach(paymentMethodID, params)

}
func (s *stripeservice) updateInvoiceSettingsToDefault(pm *stripe.PaymentMethod, nc *stripe.Customer) error {
	customerParams := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(pm.ID),
		},
	}

	if _, err := s.Customers.Update(nc.ID, customerParams); err != nil {
		return err
	}

	return nil
}
