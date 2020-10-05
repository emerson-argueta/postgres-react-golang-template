package administrator

import (
	"emersonargueta/m/v1/domain"
)

const role = "administrator"

// Churches represents a map of churches managed by an administrator where the
// key is a churchID and the value conatins the administrator's church.
type Churches map[int64]*Church

// Keys represent the church id
func (c *Churches) Keys() []int64 {
	keys := make([]int64, len(*c))

	i := 0
	for k := range *c {
		keys[i] = k
		i++
	}
	return keys
}

// Church represents a church managed by an administrator. An
// administrator has an access level and role within a church.
type Church struct {
	Access *domain.Access `json:"access,omitempty"`
	Role   *domain.Role   `json:"role,omitempty"`
}

// SubscriptionType can be Basic,Standard or Premium
type SubscriptionType string

// SubscriptionType constants
const (
	FreePlan     = SubscriptionType("free")
	BasicPlan    = SubscriptionType("basic")
	StandardPlan = SubscriptionType("standard")
	PremiumPlan  = SubscriptionType("premium")
)

func (s SubscriptionType) String() string {
	return string(s)
}

// Subscription reperesents the service plan that an administrator is currently
// registered in. The type can be free,basic,standard,premium,etc. The paymentgateway
// can represent necessary data to integrate with payment gateways such as
// stripe,paypal,square,etc.
type Subscription struct {
	Freeusagelimitcount *int64                 `json:"freeusagelimitcount,omitempty"`
	Customeremail       *string                `json:"customeremail,omitempty"`
	Type                *SubscriptionType      `json:"type,omitempty"`
	Paymentgateway      *domain.PaymentGateway `json:"paymentgateway,omitempty"`
}

// Administrator represents a user who manages churches,donators,donations, and
// possibly other administrators.
type Administrator struct {
	UUID         *string       `db:"uuid" json:"uuid"`
	Firstname    *string       `db:"firstname" json:"firstname"`
	Lastname     *string       `db:"lastname" json:"lastname"`
	Address      *string       `db:"address" json:"address"`
	Phone        *string       `db:"phone" json:"phone"`
	Churches     *Churches     `db:"churches" json:"churches,omitempty"`
	Subscription *Subscription `db:"subscription" json:"subscription,omitempty"`
	Email        *string       `json:"email,omitempty"`
	Password     *string       `json:"password,omitempty"`
}

// Client creates a connection to a service. In this case, an administrator service.
type Client interface {
	Service() Service
}

// Service provides functions to manage an
// administrator.
type Service interface {
	CreateManagementSession() error
	EndManagementSession() error
	Create(*Administrator) error
	Read(*Administrator) (*Administrator, error)
	ReadMultiple(administratorUUIDs []string) ([]*Administrator, error)
	Update(*Administrator) error
	Delete(*Administrator) error
}
