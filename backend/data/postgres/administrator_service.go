package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	"trustdonations.org/m/v2/domain/administrator"
)

const administratorTable = "administrator"

var _ administrator.Service = &Administrator{}

// Administrator represents a service for managing an administrator.
type Administrator struct {
	client *Client
}

// CreateManagementSession opens a session to start a series of actions taken to
// manage an administrator. Opening a session makes it possible to run multiple
// queries and rollback if any of the queries fail or commit if all queries are
// successful.
func (s *Administrator) CreateManagementSession() error {
	return s.client.createManagementSession()
}

// EndManagementSession ends the session created to execute a series of actions
// taken to manage an administrator. Ending a session makes it possible to
// rollback if any of the queries in a management session fail or commit if all
// queries are successful.
func (s *Administrator) EndManagementSession() error {
	return s.client.endManagementSession()
}

// Create a new administrator.
func (s *Administrator) Create(a *administrator.Administrator) (e error) {
	schema := s.client.config.Database.Schema

	query, e := NewQuery(a)
	if e != nil {
		return e
	}

	administratorInsertQuery := query.Create(schema, administratorTable)
	administratorInsertQuery = s.client.db.Rebind(administratorInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(a, administratorInsertQuery, queryParams...)
	} else {
		e = s.client.db.Get(a, administratorInsertQuery, queryParams...)
	}

	var uniqueViolation pq.ErrorCode = "23505"
	if pqError, ok := e.(*pq.Error); e != nil && !ok {
		return e
	} else if pqError != nil && pqError.Code == uniqueViolation {
		return administrator.ErrAdministratorExists
	} else if pqError != nil {
		return pqError
	}

	return nil
}

// Read an administrator by uuid.
func (s *Administrator) Read(a *administrator.Administrator) (res *administrator.Administrator, e error) {
	schema := s.client.config.Database.Schema

	filter := "UUID=?"
	queryParam := a.UUID

	query, err := NewQuery(a)
	if err != nil {
		return nil, err
	}

	administratorSelectQuery := query.Read(schema, administratorTable, filter)
	administratorSelectQuery = s.client.db.Rebind(administratorSelectQuery)

	res = &administrator.Administrator{}
	// if err := s.client.transaction.Get(res, administratorSelectQuery, queryParam); err == sql.ErrNoRows {
	if err := s.client.db.Get(res, administratorSelectQuery, queryParam); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil

}

//ReadMultiple administrators given an array of administrator ids
func (s *Administrator) ReadMultiple(administratorUUIDs []string) (res []*administrator.Administrator, e error) {
	schema := s.client.config.Database.Schema

	query, err := NewQuery(&administrator.Administrator{})
	if err != nil {
		return nil, err
	}
	var queryParams []interface{}
	for _, elem := range administratorUUIDs {
		queryParams = append(queryParams, elem)
	}

	filter := query.CreateMultipleValueFilter("UUID", len(queryParams))

	administratorSelectQuery := query.Read(schema, administratorTable, filter)
	administratorSelectQuery = s.client.db.Rebind(administratorSelectQuery)

	// if err := s.client.transaction.Select(&res, administratorSelectQuery, queryParams...); err == sql.ErrNoRows {
	if err := s.client.db.Select(&res, administratorSelectQuery, queryParams...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

// Update an administrator, searching by uuid.
func (s *Administrator) Update(a *administrator.Administrator) (e error) {
	schema := s.client.config.Database.Schema
	filter := "UUID=?"

	query, e := NewQuery(a)
	if e != nil {
		return e
	}

	administratorUpdateQuery := query.Update(schema, administratorTable, filter)
	administratorUpdateQuery = s.client.db.Rebind(administratorUpdateQuery)

	includeNil := true
	queryParams := append(query.ModelValues(includeNil), a.UUID)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(a, administratorUpdateQuery, queryParams...)
	} else {
		e = s.client.db.Get(a, administratorUpdateQuery, queryParams...)
	}
	if e != nil {
		return e
	}

	return nil
}

// Delete an adminstrator searching by uuid.
func (s *Administrator) Delete(a *administrator.Administrator) (e error) {
	schema := s.client.config.Database.Schema

	filter := "UUID=?"
	queryParam := a.UUID

	query, e := NewQuery(a)
	if e != nil {
		return e
	}
	administratorDeleteQuery := query.Delete(schema, administratorTable, filter)
	administratorDeleteQuery = s.client.db.Rebind(administratorDeleteQuery)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(a, administratorDeleteQuery, queryParam)
	} else {
		e = s.client.db.Get(a, administratorDeleteQuery, queryParam)
	}

	if e != nil && e == sql.ErrNoRows {
		return administrator.ErrAdministratorNotFound
	} else if e != nil {
		return e
	}

	return nil
}
