package postgres

import (
	"database/sql"
	"emersonargueta/m/v1/communitygoaltracker"
	"emersonargueta/m/v1/communitygoaltracker/achiever"

	"github.com/lib/pq"
)

var _ achiever.Processes = &achieverservice{}

// achieverservice represents a service for managing a Achiever.
type achieverservice struct {
	client *Client
}

// CreateAchiever if successful. If the achiever
// exists, returns ErrAchieverExists.
func (s *achieverservice) CreateAchiever(a *achiever.Achiever) (res *achiever.Achiever, e error) {
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

// RetrieveAchiever by uuid. If the achiever does not exists,
// returns ErrAchieverNotFound.
func (s *achieverservice) RetrieveAchiever(uuid string) (res *achiever.Achiever, e error) {
	filter := "UUID=?"
	queryParam := uuid

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

// RetrieveAchievers searching by uuids. Returns ErrAchiverNotFound if none of the achievers
// are found.
func (s *achieverservice) RetrieveAchievers(uuids []string) (res []*achiever.Achiever, e error) {

	query, err := NewQuery(&achiever.Achiever{})
	if err != nil {
		return nil, err
	}
	var queryParams []interface{}
	for _, elem := range uuids {
		queryParams = append(queryParams, elem)
	}

	filter := query.CreateMultipleValueFilter("UUID", len(queryParams))

	achieverSelectQuery := query.Read(CommunitygoaltrackerSchema, AchieverTable, filter)
	achieverSelectQuery = s.client.db.Rebind(achieverSelectQuery)

	if err := s.client.db.Select(&res, achieverSelectQuery, queryParams...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateAchiever searching by uuid. If the achiever does not exists, returns
// ErrAchieverNotFound.
func (s *achieverservice) UpdateAchiever(a *achiever.Achiever) (e error) {
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
func (s *achieverservice) DeleteAchiever(uuid string) (e error) {
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
