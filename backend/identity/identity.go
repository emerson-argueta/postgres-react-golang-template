package identity

import (
	"emersonargueta/m/v1/identity/domain"
	"emersonargueta/m/v1/identity/user"
)

// DomainName of this package
const DomainName = "identity"

// Identity exposes identity domain processes.
type Identity struct {
	client *Client
	// Services from domain models used internally in identity domain processes
	User   user.Service
	Domain domain.Service
}

// RegisterUser using the following business logic
func (i *Identity) RegisterUser(u *user.User) (res *user.User, e error) {

	return res, e
}

//LoginUser using the following business logic
func (i *Identity) LoginUser(email string, password string) (res *user.User, e error) {
	return res, e
}

//UpdateUser using the following business logic
func (i *Identity) UpdateUser(u *user.User) (e error) {
	return e
}

//UnRegisterUser using the following business logic
func (i *Identity) UnRegisterUser(u *user.User) (e error) {
	return e
}

//AddDomain using the following business logic
func (i *Identity) AddDomain(d *domain.Domain) (e error) {
	return e
}

//LookupDomain using the following business logic
func (i *Identity) LookupDomain(d *domain.Domain) (res *domain.Domain, e error) {
	return res, e
}

//UpdateDomain using the following business logic
func (i *Identity) UpdateDomain(d *domain.Domain) (e error) {
	return e
}

//RemoveDomain using the following business logic
func (i *Identity) RemoveDomain(d *domain.Domain) (e error) {
	return e
}
