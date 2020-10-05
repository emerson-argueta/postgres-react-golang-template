package donator

import (
	"time"

	"emersonargueta/m/v1/domain"
)

// Churches represents a map of churches that a donator donates to where the key
// is the church ID and the value is a church.
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

// Church represents a church that a donator donates to. The donator contains a
// count of donations and the time of first donation for that church.
type Church struct {
	Donationcount *int64     `json:"donationcount,omitempty"`
	Firstdonation *time.Time `json:"firstdonation,omitempty"`
}

// Donator represents a donator that donates to a church who is managed by an administrator.
type Donator struct {
	ID               *int64                   `db:"id" dbignoreinsert:"" json:"id"`
	UUID             *string                  `db:"uuid" json:"uuid"`
	Firstname        *string                  `db:"firstname" json:"firstname"`
	Lastname         *string                  `db:"lastname" json:"lastname"`
	Email            *string                  `db:"email" json:"email"`
	Address          *string                  `db:"address" json:"address"`
	Phone            *string                  `db:"phone" json:"phone"`
	Churches         *Churches                `db:"churches" json:"churches"`
	Accountstatement *domain.AccountStatement `db:"accountstatement" json:"accountstatement"`
}

// Client creates a connection to a service. In this case, a donator service.
type Client interface {
	Service() Service
}

// Service provides functions that can be used to manage donators.
type Service interface {
	CreateManagementSession() error
	EndManagementSession() error
	Create(*Donator) error
	Read(d *Donator, byEmail bool) (*Donator, error)
	ReadMultiple(donatorids []int64) (res []*Donator, e error)
	// ReadWithFilter a donator, searching by non empty fields of filterDonator model.
	ReadWithFilter(d *Donator, filterDonator *Donator) (*Donator, error)
	Update(d *Donator, byEmail bool) error
	Delete(d *Donator, byEmail bool) error
}
