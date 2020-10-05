package stripe

import (
	"emersonargueta/m/v1/domain/administrator"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
)

var _ administrator.SubscriptionActions = &Administrator{}

// Administrator represents a service for managing stripe subscriptions.
type Administrator struct {
	client  *Client
	Usecase administrator.Usecase
}

// NewSubscription for an administrator with stripe. A new subscription for an
// administrator can be created so that the administrator can gain access to the
// church_fund_managing service after reaching the free usage limit
// (100 automated donation entries and 400 manual donation entries).
func (s *Administrator) NewSubscription(sb *administrator.Subscription, a *administrator.Administrator) error {
	newCustomer, err := createCustomer(s, sb)
	if err != nil {
		return ErrStripeNewSubscription
	}

	newSubscription, err := createSubscription(s, sb, newCustomer)

	subscriptionID := newSubscription.ID
	subscriptionEnd := newSubscription.CurrentPeriodEnd
	subscription := map[string]interface{}{
		"subscriptionID":  subscriptionID,
		"subscriptionEnd": subscriptionEnd,
	}

	customerID := newSubscription.Customer.ID
	customerEmail := newSubscription.Customer.Email
	customer := map[string]interface{}{
		"customerID":    customerID,
		"customerEmail": customerEmail,
	}

	planID := newSubscription.Items.Data[0].Plan.ID
	planName := newSubscription.Items.Data[0].Plan.Nickname
	planAmount := newSubscription.Items.Data[0].Plan.Amount
	planProductID := newSubscription.Items.Data[0].Plan.Product.ID
	plan := map[string]interface{}{
		"planID":        planID,
		"planName":      planName,
		"planAmount":    planAmount,
		"planProductID": planProductID,
	}

	paymentGateway := *sb.Paymentgateway
	paymentGateway["subscription"] = subscription
	paymentGateway["customer"] = customer
	paymentGateway["plan"] = plan

	// Execute business logic to finalize the new subscription.
	s.Usecase.NewSubscription(sb, a)

	return nil
}
func createCustomer(s *Administrator, sb *administrator.Subscription) (*stripe.Customer, error) {
	stripeEmail := sb.Customeremail
	customerParams := stripe.CustomerParams{Email: stripeEmail}
	newCustomer, err := s.client.stripe.Customers.New(&customerParams)
	if err != nil {
		return nil, err
	}

	return newCustomer, nil
}
func createSubscription(s *Administrator, sb *administrator.Subscription, newCustomer *stripe.Customer) (*stripe.Subscription, error) {
	paymentMethod, err := attachPaymentMethod(s, sb, newCustomer)
	if err != nil {
		return nil, ErrStripeNewSubscription
	}

	if err := updateInvoiceSettingsToDefault(s, paymentMethod, newCustomer); err != nil {
		return nil, ErrStripeNewSubscription
	}

	subscriptionParams := &stripe.SubscriptionParams{
		Customer: &newCustomer.ID,
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(sb.Type.String()),
			},
		},
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")

	return sub.New(subscriptionParams)
}
func attachPaymentMethod(s *Administrator, sb *administrator.Subscription, newCustomer *stripe.Customer) (*stripe.PaymentMethod, error) {
	paymentGateway := *sb.Paymentgateway
	paymentMethodID, ok := paymentGateway["paymentMethodID"].(string)
	if !ok {
		return nil, administrator.ErrAdministratorSubscriptionPaymentGateway
	}

	params := &stripe.PaymentMethodAttachParams{
		Customer: &newCustomer.ID,
	}

	return s.client.stripe.PaymentMethods.Attach(paymentMethodID, params)

}
func updateInvoiceSettingsToDefault(s *Administrator, pm *stripe.PaymentMethod, nc *stripe.Customer) error {
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
