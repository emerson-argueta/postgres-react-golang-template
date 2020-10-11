package postgres

import (
	"database/sql"
	"emersonargueta/m/v1/communitygoaltracker"
	"emersonargueta/m/v1/communitygoaltracker/achiever"

	"github.com/lib/pq"
)

var _ achiever.Service = &Achiever{}

// Achiever represents a service for managing a Achiever.
type Achiever struct {
	client *Client
}

// CreateAchiever if successful. If the achiever
// exists, returns ErrAchieverExists.
func (s *Achiever) CreateAchiever(a *achiever.Achiever) (res *achiever.Achiever, e error) {
	query, e := NewQuery(a)
	if e != nil {
		return nil, e
	}

	achieverInsertQuery := query.Create(CommunitygoaltrackerSchema, AchieverTable)
	achieverInsertQuery = s.client.db.Rebind(achieverInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	res = &achiever.Achiever{}

	e = s.client.db.Get(res, achieverInsertQuery, queryParams...)

	var uniqueViolation pq.ErrorCode = "23505"
	if pqError, ok := e.(*pq.Error); e != nil && !ok {
		return nil, e
	} else if pqError != nil && pqError.Code == uniqueViolation {
		return nil, communitygoaltracker.ErrAchieverExists
	} else if pqError != nil {
		return nil, pqError
	}

	return res, e
}

// RetrieveAchiever by email. If the achiever does not exists,
// returns ErrAchieverNotFound.
func (s *Achiever) RetrieveAchiever(email string) (res *achiever.Achiever, e error) {
	filter := "EMAIL=?"
	queryParam := email

	query, e := NewQuery(&achiever.Achiever{})
	if e != nil {
		return nil, e
	}

	achieverSelectQuery := query.Read(CommunitygoaltrackerSchema, AchieverTable, filter)
	achieverSelectQuery = s.client.db.Rebind(achieverSelectQuery)

	res = &achiever.Achiever{}
	e = s.client.db.Get(res, achieverSelectQuery, queryParam)
	if e == sql.ErrNoRows {
		return nil, communitygoaltracker.ErrAchieverNotFound
	}

	return res, e

}

// UpdateAchiever searching by uuid. If the achiever does not exists, returns
// ErrAchieverNotFound.
func (s *Achiever) UpdateAchiever(a *achiever.Achiever) (e error) {
	filter := "UUID=?"
	queryParam := a.UUID

	query, e := NewQuery(a)
	if e != nil {
		return e
	}

	achieverUpdateQuery := query.Update(CommunitygoaltrackerSchema, AchieverTable, filter)
	achieverUpdateQuery = s.client.db.Rebind(achieverUpdateQuery)

	includeNil := true
	queryParams := append(query.ModelValues(includeNil), queryParam)

	e = s.client.db.Get(a, achieverUpdateQuery, queryParams...)

	if e == sql.ErrNoRows {
		return communitygoaltracker.ErrAchieverNotFound
	}

	return e
}

// DeleteAchiever searching by uuid. If the achiever does not exists, returns
// ErrAchieverNotFound.
func (s *Achiever) DeleteAchiever(uuid string) (e error) {
	filter := "UUID=?"
	queryParam := uuid

	query, e := NewQuery(&achiever.Achiever{})
	if e != nil {
		return e
	}
	achieverDeleteQuery := query.Delete(CommunitygoaltrackerSchema, AchieverTable, filter)
	achieverDeleteQuery = s.client.db.Rebind(achieverDeleteQuery)

	e = s.client.db.Get(&achiever.Achiever{}, achieverDeleteQuery, queryParam)
	if e == sql.ErrNoRows {
		return communitygoaltracker.ErrAchieverNotFound
	}

	return e
}
