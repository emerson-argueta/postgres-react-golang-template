package stripe

import (
	"emersonargueta/m/v1/paymentgateway"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
)

var _ paymentgateway.Service = &Stripe{}

// Stripe represents a service for managing stripe subscriptions.
type Stripe struct {
	client *Client
}

// NewSubscription for a user with stripe. A new subscription for a
// user can be created so that the user can gain access to
// a domain.
func (s *Stripe) NewSubscription(email string, planType string, paymentMethodID string) (re *stripe.Subscription, e error) {
	newCustomer, err := s.createCustomer(email)
	if err != nil {
		return nil, ErrStripeNewSubscription
	}

	newSubscription, err := s.createSubscription(planType, paymentMethodID, newCustomer)
	if err != nil {
		return nil, err
	}

	// subscriptionID := newSubscription.ID
	// subscriptionEnd := newSubscription.CurrentPeriodEnd
	// subscription := map[string]interface{}{
	// 	"subscriptionID":  subscriptionID,
	// 	"subscriptionEnd": subscriptionEnd,
	// }

	// customerID := newSubscription.Customer.ID
	// customerEmail := newSubscription.Customer.Email
	// customer := map[string]interface{}{
	// 	"customerID":    customerID,
	// 	"customerEmail": customerEmail,
	// }

	// planID := newSubscription.Items.Data[0].Plan.ID
	// planName := newSubscription.Items.Data[0].Plan.Nickname
	// planAmount := newSubscription.Items.Data[0].Plan.Amount
	// planProductID := newSubscription.Items.Data[0].Plan.Product.ID
	// plan := map[string]interface{}{
	// 	"planID":        planID,
	// 	"planName":      planName,
	// 	"planAmount":    planAmount,
	// 	"planProductID": planProductID,
	// }

	// paymentGateway := *sb.Paymentgateway
	// paymentGateway["subscription"] = subscription
	// paymentGateway["customer"] = customer
	// paymentGateway["plan"] = plan

	// // Execute business logic to finalize the new subscription.
	// s.Usecase.NewSubscription(sb, a)

	return newSubscription, nil
}
func (s *Stripe) createCustomer(email string) (*stripe.Customer, error) {
	customerParams := stripe.CustomerParams{Email: &email}
	newCustomer, err := s.client.stripe.Customers.New(&customerParams)
	if err != nil {
		return nil, err
	}

	return newCustomer, nil
}
func (s *Stripe) createSubscription(planType string, paymentMethodID string, newCustomer *stripe.Customer) (*stripe.Subscription, error) {
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
func (s *Stripe) attachPaymentMethod(paymentMethodID string, newCustomer *stripe.Customer) (*stripe.PaymentMethod, error) {

	params := &stripe.PaymentMethodAttachParams{
		Customer: &newCustomer.ID,
	}

	return s.client.stripe.PaymentMethods.Attach(paymentMethodID, params)

}
func (s *Stripe) updateInvoiceSettingsToDefault(pm *stripe.PaymentMethod, nc *stripe.Customer) error {
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
