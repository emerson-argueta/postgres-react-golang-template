package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	"trustdonations.org/m/v2/domain/donator"
)

const donatorTable = "donator"

var _ donator.Service = &Donator{}

// Donator represents a service for managing a donator.
type Donator struct {
	client *Client
}

// CreateManagementSession opens a session to start a series of actions taken to
// manage a donator. Opening a session makes it possible to run multiple
// queries and rollback if any of the queries fail or commit if all queries are
// successful.
func (s *Donator) CreateManagementSession() error {
	return s.client.createManagementSession()
}

// EndManagementSession ends the session created to execute a series of actions
// taken to manage a donator. Ending a session makes it possible to
// rollback if any of the queries in a management session fail or commit if all
// queries are successful.
func (s *Donator) EndManagementSession() error {
	return s.client.endManagementSession()
}

// Create a new donator.
func (s *Donator) Create(d *donator.Donator) (e error) {
	schema := s.client.config.Database.Schema

	query, e := NewQuery(d)
	if e != nil {
		return e
	}

	donatorInsertQuery := query.Create(schema, donatorTable)
	donatorInsertQuery = s.client.db.Rebind(donatorInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(d, donatorInsertQuery, queryParams...)
	} else {
		e = s.client.db.Get(d, donatorInsertQuery, queryParams...)
	}

	var uniqueViolation pq.ErrorCode = "23505"
	if pqError, ok := e.(*pq.Error); e != nil && !ok {
		return e
	} else if pqError != nil && pqError.Code == uniqueViolation {
		return donator.ErrDonatorExists
	} else if pqError != nil {
		return pqError
	}

	return nil
}

// Read a donator, searching by id or email.
func (s *Donator) Read(d *donator.Donator, byEmail bool) (res *donator.Donator, e error) {
	schema := s.client.config.Database.Schema

	filter := "ID=?"
	var queryParam interface{} = d.ID
	if byEmail {
		filter = "EMAIL=?"
		queryParam = d.Email
	}

	query, err := NewQuery(d)
	if err != nil {
		return nil, err
	}

	donatorSelectQuery := query.Read(schema, donatorTable, filter)
	donatorSelectQuery = s.client.db.Rebind(donatorSelectQuery)

	res = &donator.Donator{}
	if err := s.client.db.Get(res, donatorSelectQuery, queryParam); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil

}

// ReadMultiple donators given an array of ids
func (s *Donator) ReadMultiple(donatorids []int64) (res []*donator.Donator, e error) {
	schema := s.client.config.Database.Schema

	query, err := NewQuery(&donator.Donator{})
	if err != nil {
		return nil, err
	}
	var queryParams []interface{}
	for _, elem := range donatorids {
		queryParams = append(queryParams, elem)
	}

	filter := query.CreateMultipleValueFilter("ID", len(queryParams))

	donatorSelectQuery := query.Read(schema, donatorTable, filter)
	donatorSelectQuery = s.client.db.Rebind(donatorSelectQuery)

	if err := s.client.db.Select(&res, donatorSelectQuery, queryParams...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil

}

// ReadWithFilter a donator, searching by non empty fields of filterDonator model.
func (s *Donator) ReadWithFilter(d *donator.Donator, filterDonator *donator.Donator) (res *donator.Donator, e error) {
	schema := s.client.config.Database.Schema

	queryFilter, err := NewQuery(filterDonator)
	if err != nil {
		return nil, err
	}
	filter := queryFilter.CreateFilterFromModel()
	includeNil := false
	queryParams := queryFilter.ModelValues(includeNil)

	query, err := NewQuery(d)
	if err != nil {
		return nil, err
	}

	donatorSelectQuery := query.Read(schema, donatorTable, filter)
	donatorSelectQuery = s.client.db.Rebind(donatorSelectQuery)

	res = &donator.Donator{}
	if err := s.client.db.Get(res, donatorSelectQuery, queryParams...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

// Update donator, searching by id or email.
func (s *Donator) Update(d *donator.Donator, byEmail bool) (e error) {
	schema := s.client.config.Database.Schema

	filter := "ID=?"
	var queryParam interface{} = d.ID
	if byEmail {
		filter = "EMAIL=?"
		queryParam = d.Email
	}

	query, e := NewQuery(d)
	if e != nil {
		return e
	}

	donatorUpdateQuery := query.Update(schema, donatorTable, filter)
	donatorUpdateQuery = s.client.db.Rebind(donatorUpdateQuery)

	includeNil := true
	queryParams := append(query.ModelValues(includeNil), queryParam)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(d, donatorUpdateQuery, queryParams...)
	} else {
		e = s.client.db.Get(d, donatorUpdateQuery, queryParams...)
	}
	if e != nil {
		return e
	}

	return nil
}

// Delete a donator searching by id or email.
func (s *Donator) Delete(d *donator.Donator, byEmail bool) (e error) {
	schema := s.client.config.Database.Schema

	filter := "ID=?"
	var queryParam interface{} = d.ID
	if byEmail {
		filter = "EMAIL=?"
		queryParam = d.Email
	}

	query, e := NewQuery(d)
	if e != nil {
		return e
	}
	donatorDeleteQuery := query.Delete(schema, donatorTable, filter)
	donatorDeleteQuery = s.client.db.Rebind(donatorDeleteQuery)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(d, donatorDeleteQuery, queryParam)
	} else {
		e = s.client.db.Get(d, donatorDeleteQuery, queryParam)
	}
	if e != nil && e == sql.ErrNoRows {
		return donator.ErrDonatorNotFound
	} else if e != nil {
		return e
	}

	return nil
}
