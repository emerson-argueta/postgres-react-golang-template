package church

import (
	"emersonargueta/m/v1/domain"
)

// Administrators represensts a map of administrators within a church where the
// key is an administrator UUID and the value contains the adminstrator.
type Administrators map[string]*Administrator

// Keys represent the administrator uuid
func (c *Administrators) Keys() []string {
	keys := make([]string, len(*c))

	i := 0
	for k := range *c {
		keys[i] = k
		i++
	}
	return keys
}

// Administrator represents an administrator with a church with a role and
// access level.
type Administrator struct {
	Access *domain.Access `json:"access,omitempty"`
	Role   *domain.Role   `json:"role,omitempty"`
}

// Donators represents a map of donators within a church where the key is a
// donator ID and the value contains the donator.
type Donators map[int64]*Donator

// Keys represent the donator id
func (c *Donators) Keys() []int64 {
	keys := make([]int64, len(*c))

	i := 0
	for k := range *c {
		keys[i] = k
		i++
	}
	return keys
}

// Donator represents a donator wihtin a church where UUID is not an empty
// string if the donator has a user account.
type Donator struct {
	UUID *string `json:"uuid,omitempty"`
}

// DonationCategories represents the categories that donations can have for a given church
type DonationCategories map[string]*string

// Keys represent the donation category
func (c *DonationCategories) Keys() []string {
	keys := make([]string, len(*c))

	i := 0
	for k := range *c {
		keys[i] = k
		i++
	}
	return keys
}

// Church represents a church that is managed by administrators and recieves
// donations from donators.
type Church struct {
	ID                 *int64                   `db:"id" dbignoreinsert:"" json:"id"`
	Type               *string                  `db:"type" json:"type"`
	Name               *string                  `db:"name" json:"name"`
	Address            *string                  `db:"address" json:"address"`
	Phone              *string                  `db:"phone" json:"phone"`
	Email              *string                  `db:"email" json:"email"`
	Password           *string                  `db:"password" json:"password,omitempty"`
	Administrators     *Administrators          `db:"administrators" json:"administrators"`
	Donators           *Donators                `db:"donators" json:"donators"`
	Accountstatement   *domain.AccountStatement `db:"accountstatement" json:"accountstatement"`
	DonationCategories *DonationCategories      `db:"donationcategories" json:"donationcategories"`
}

// Client creates a connection to a service. In this case, an church service.
type Client interface {
	Service() Service
}

// Service provides functions that can be used for managing churches.
type Service interface {
	CreateManagementSession() error
	EndManagementSession() error
	Create(*Church) error
	// Read a church byEmail or Id if byEmail flag set to false
	Read(c *Church, byEmail bool) (*Church, error)
	// Update a church ,finding church byEmail or Id if byEmail flag set to false
	Update(c *Church, byEmail bool) error
	// Delete a church ,finding church byEmail or Id if byEmail flag set to false
	Delete(c *Church, byEmail bool) error
	//ReadMultiple churches given an array of church ids
	ReadMultiple(churchids []int64) ([]*Church, error)
}
