package paymentgateway

// Processes used for a payment gateway.
type Processes interface {
	// NewPayment made by customer.
	NewPayment(details map[string]string) (e error)
	// NewSubscription for a customer to gain access to services.
	NewSubscription(details map[string]string) (e error)
}
