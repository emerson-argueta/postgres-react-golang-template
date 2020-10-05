package postgres

import (
	"database/sql"

	"emersonargueta/m/v1/domain/church"

	"github.com/lib/pq"
)

const churchTable = "church"

var _ church.Service = &Church{}

// Church represents a service for managing a church.
type Church struct {
	client *Client
}

// CreateManagementSession opens a session to start a series of actions taken to
// manage an church. Opening a session makes it possible to run multiple
// queries and rollback if any of the queries fail or commit if all queries are
// successful.
func (s *Church) CreateManagementSession() error {
	return s.client.createManagementSession()
}

// EndManagementSession ends the session created to execute a series of actions
// taken to manage an church. Ending a session makes it possible to
// rollback if any of the queries in a management session fail or commit if all
// queries are successful.
func (s *Church) EndManagementSession() error {
	return s.client.endManagementSession()
}

// Create a new church.
func (s *Church) Create(c *church.Church) (e error) {
	schema := s.client.config.Database.Schema

	query, e := NewQuery(c)
	if e != nil {
		return e
	}

	churchInsertQuery := query.Create(schema, churchTable)
	churchInsertQuery = s.client.db.Rebind(churchInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(c, churchInsertQuery, queryParams...)
	} else {
		e = s.client.db.Get(c, churchInsertQuery, queryParams...)
	}
	var uniqueViolation pq.ErrorCode = "23505"
	if pqError, ok := e.(*pq.Error); e != nil && !ok {
		return e
	} else if pqError != nil && pqError.Code == uniqueViolation {
		return church.ErrChurchExists
	} else if pqError != nil {
		return pqError
	}

	return nil
}

// Read a church, searching by id or email.
func (s *Church) Read(c *church.Church, byEmail bool) (res *church.Church, e error) {
	schema := s.client.config.Database.Schema

	filter := "ID=?"
	var queryParam interface{} = c.ID
	if byEmail {
		filter = "EMAIL=?"
		queryParam = c.Email
	}

	query, err := NewQuery(c)
	if err != nil {
		return nil, err
	}

	churchSelectQuery := query.Read(schema, churchTable, filter)
	churchSelectQuery = s.client.db.Rebind(churchSelectQuery)

	res = &church.Church{}
	// if err := s.client.transaction.Get(res, churchSelectQuery, queryParam); err == sql.ErrNoRows {
	if err := s.client.db.Get(res, churchSelectQuery, queryParam); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil

}

//ReadMultiple churches given an array of church ids
func (s *Church) ReadMultiple(churchids []int64) (res []*church.Church, e error) {
	schema := s.client.config.Database.Schema

	query, err := NewQuery(&church.Church{})
	if err != nil {
		return nil, err
	}
	var queryParams []interface{}
	for _, elem := range churchids {
		queryParams = append(queryParams, elem)
	}

	filter := query.CreateMultipleValueFilter("ID", len(queryParams))

	churchSelectQuery := query.Read(schema, churchTable, filter)
	churchSelectQuery = s.client.db.Rebind(churchSelectQuery)

	// if err := s.client.transaction.Select(&res, churchSelectQuery, queryParams...); err == sql.ErrNoRows {
	if err := s.client.db.Select(&res, churchSelectQuery, queryParams...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil

}

// Update the church, searching by id or email.
func (s *Church) Update(c *church.Church, byEmail bool) (e error) {
	schema := s.client.config.Database.Schema

	filter := "ID=?"
	var queryParam interface{} = c.ID
	if byEmail {
		filter = "EMAIL=?"
		queryParam = c.Email
	}

	query, e := NewQuery(c)
	if e != nil {
		return e
	}

	churchUpdateQuery := query.Update(schema, churchTable, filter)
	churchUpdateQuery = s.client.db.Rebind(churchUpdateQuery)

	includeNil := true
	queryParams := append(query.ModelValues(includeNil), queryParam)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(c, churchUpdateQuery, queryParams...)
	} else {
		e = s.client.db.Get(c, churchUpdateQuery, queryParams...)
	}
	if e != nil {
		return e
	}

	return nil
}

// Delete a church searching by id or email.
func (s *Church) Delete(c *church.Church, byEmail bool) (e error) {
	schema := s.client.config.Database.Schema

	filter := "ID=?"
	var queryParam interface{} = c.ID
	if byEmail {
		filter = "EMAIL=?"
		queryParam = c.Email
	}

	query, e := NewQuery(c)
	if e != nil {
		return e
	}
	churchDeleteQuery := query.Delete(schema, churchTable, filter)
	churchDeleteQuery = s.client.db.Rebind(churchDeleteQuery)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(c, churchDeleteQuery, queryParam)
	} else {
		e = s.client.db.Get(c, churchDeleteQuery, queryParam)
	}
	if e != nil && e == sql.ErrNoRows {
		return church.ErrChurchNotFound
	} else if e != nil {
		return e
	}

	return nil
}
