package identity

import (
	"emersonargueta/m/v1/identity/domain"
	"emersonargueta/m/v1/identity/user"
)

// Client creates a connection to a service. In this case, an identity service.
type Client interface {
	Service() Service
}

// Service provides processes provided by the identity domain.
type Service interface {
	user.Service
	domain.Service
}

// Usecase for process logic.
type Usecase struct {
	Services Services
}

// Services used by usecase
type Services struct {
	User   user.Service
	Domain domain.Service
}

// RegisterUser using the following business logic
func (uc *Usecase) RegisterUser(u *user.User) (e error) {

	return e
}

//LoginUser using the following business logic
func (uc *Usecase) LoginUser(email string, password string) (res *user.User, e error) {
	return res, e
}

//UpdateUser using the following business logic
func (uc *Usecase) UpdateUser(u *user.User) (e error) {
	return e
}

//UnRegisterUser using the following business logic
func (uc *Usecase) UnRegisterUser(u *user.User) (e error) {
	return e
}

//AddDomain using the following business logic
func (uc *Usecase) AddDomain(d *domain.Domain) (e error) {
	return e
}

//LookupDomain using the following business logic
func (uc *Usecase) LookupDomain(d *domain.Domain) (res *domain.Domain, e error) {
	return res, e
}

//UpdateDomain using the following business logic
func (uc *Usecase) UpdateDomain(d *domain.Domain) (e error) {
	return e
}

//RemoveDomain using the following business logic
func (uc *Usecase) RemoveDomain(d *domain.Domain) (e error) {
	return e
}
