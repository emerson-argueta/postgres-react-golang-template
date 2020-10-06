package domain

//Domain model represents domains which users can be part of
type Domain struct {
	ID   *int64  `db:"id" dbignoreinsert:"" json:"id"`
	Name *string `db:"name" json:"name"`
}

// Service provides processes which affect the domains.
type Service interface {
	CreateDomain(*Domain) error
	RetrieveDomain(*Domain) (*Domain, error)
	UpdateDomain(*Domain) error
	DeleteDomain(*Domain) error
}
