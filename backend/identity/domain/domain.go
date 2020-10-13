package domain

//Domain model represents domains which users can be part of
type Domain struct {
	ID   *int64  `db:"id" dbignoreinsert:"" json:"id"`
	Name *string `db:"name" json:"name"`
}

// Processes used to modify the domain model.
type Processes interface {
	// CreateDomain implementation must return ErrDomainExists if domain exists.
	CreateDomain(*Domain) (*Domain, error)
	// RetreiveDomain implementation must return ErrDomainNotFound if domain not found.
	RetrieveDomain(name string) (*Domain, error)
	// UpdateDomain implementation must search domain by id and return
	// ErrDomainNotFound if domain not found. Must return ErrDomainExists if
	// update name conflicts with another domain.
	UpdateDomain(*Domain) error
	// DeleteDomain implementation must search domain by id and return
	// ErrDomainNotFound if domain not found.
	DeleteDomain(id int64) error
}
