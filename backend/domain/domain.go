package domain

import (
	"time"
)

// ServiceName of the application
const ServiceName = "church_fund_managing"

// Access represents the access level an administrator has within a church.
// The access level will determine how an administrator can manage other
// administrators within the same church.
type Access string

func (a Access) String() string {
	return string(a)
}

// NonRestricted enables an administrator to manage other administrators within the same church.
const NonRestricted = Access("non-restricted")

// Restricted disables an administrator to manage other administrators within the same church.
const Restricted = Access("restricted")

// Role of an administrator within the church. The role will help in determining
// how an administrator will be charged for using the church_fund_managing service.
type Role string

func (r Role) String() string {
	return string(r)
}

// Creator role within a church for an administrator. Whenever an adminstrator
// creates a church, they will have a creator role within that church. If an
// administrator is a creator of a church, every donaton made to that church
// will be counted towards their free usage limit (100 automated input
// donations, 400 manual input donations).
const Creator = Role("creator")

// Support role within a church for an administrator. Whenever an administrator
// adds a church, they will have a support role within that church. The creator
// of the church will be charged for allowing administrators with support roles
// into their church.
const Support = Role("support")

// AccountStatement represents a statement containing the date and monthly
// transaction sum or closing balance for a church or donator.
type AccountStatement struct {
	Closingbalance *float64   `json:"closingbalance,omitempty"`
	Date           *time.Time `json:"date,omitempty"`
}

// PaymentGateway contains information about whatever payment gateway(stripe,paypal,etc) is used to create subscriptions.
type PaymentGateway map[string]interface{}

// Filter is a map used to filter domain models (church, donator,
// administrator,transaction) by their field values. The fields are specified in
// the key of the map. The value or values are specified by the map's
// FilterValue struct. The operator used to filter by field and value is
// specified by the map's FilterValue struct.
type Filter map[string][]FilterValue

// FilterValue hold the value or values and comparator operator used on each value or values to filter transactions.
type FilterValue struct {
	ComparatorOperator string
	Value              interface{}
}

// Values of Filter map
func (f Filter) Values() (res []interface{}) {
	for _, value := range f {
		res = append(res, value)
	}
	return res
}
