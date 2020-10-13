package postgres

import (
	"database/sql"
	"emersonargueta/m/v1/identity"
	"emersonargueta/m/v1/identity/domain"

	"github.com/lib/pq"
)

var _ domain.Processes = &DomainService{}

// DomainService represents a service for managing a domain.
type DomainService struct {
	client *Client
}

// CreateDomain if successful. If the domain
// exists, returns ErrDomainExists.
func (s *DomainService) CreateDomain(d *domain.Domain) (res *domain.Domain, e error) {
	query, e := NewQuery(d)
	if e != nil {
		return nil, e
	}

	userInsertQuery := query.Create(IdentitySchema, DomainTable)
	userInsertQuery = s.client.db.Rebind(userInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	res = &domain.Domain{}

	e = s.client.db.Get(res, userInsertQuery, queryParams...)

	var uniqueViolation pq.ErrorCode = "23505"
	if pqError, ok := e.(*pq.Error); e != nil && !ok {
		return nil, e
	} else if pqError != nil && pqError.Code == uniqueViolation {
		return nil, identity.ErrDomainExists
	} else if pqError != nil {
		return nil, pqError
	}

	return res, e
}

// RetrieveDomain searching by name. If the domain does not exists,
// returns ErrDomainNotFound.
func (s *DomainService) RetrieveDomain(name string) (res *domain.Domain, e error) {
	filter := "NAME=?"
	queryParam := name

	query, e := NewQuery(&domain.Domain{})
	if e != nil {
		return nil, e
	}

	userSelectQuery := query.Read(IdentitySchema, DomainTable, filter)
	userSelectQuery = s.client.db.Rebind(userSelectQuery)

	res = &domain.Domain{}
	e = s.client.db.Get(res, userSelectQuery, queryParam)
	if e == sql.ErrNoRows {
		return nil, identity.ErrDomainNotFound
	}

	return res, e

}

// UpdateDomain searching by id. If the domain does not exists, returns
// ErrDomainNotFound.
func (s *DomainService) UpdateDomain(d *domain.Domain) (e error) {
	filter := "ID=?"
	queryParam := d.ID

	query, e := NewQuery(d)
	if e != nil {
		return e
	}

	userUpdateQuery := query.Update(IdentitySchema, DomainTable, filter)
	userUpdateQuery = s.client.db.Rebind(userUpdateQuery)

	includeNil := true
	queryParams := append(query.ModelValues(includeNil), queryParam)

	e = s.client.db.Get(d, userUpdateQuery, queryParams...)
	if e == sql.ErrNoRows {
		return identity.ErrDomainNotFound
	}

	return e
}

// DeleteDomain searching by id. If the domain does not exists, returns
// ErrDomainNotFound.
func (s *DomainService) DeleteDomain(id int64) (e error) {
	filter := "ID=?"
	queryParam := id

	query, e := NewQuery(&domain.Domain{})
	if e != nil {
		return e
	}
	userDeleteQuery := query.Delete(IdentitySchema, DomainTable, filter)
	userDeleteQuery = s.client.db.Rebind(userDeleteQuery)

	e = s.client.db.Get(&domain.Domain{}, userDeleteQuery, queryParam)
	if e == sql.ErrNoRows {
		return identity.ErrDomainNotFound
	}

	return e
}
