package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	"trustdonations.org/m/v2/user"
)

const userTable = "user"
const schema = "identity"

var _ user.Service = &User{}

// User represents a service for managing a user.
type User struct {
	client *Client
}

// CreateManagementSession opens a session to start a series of actions taken to
// manage an user. Opening a session makes it possible to run multiple
// queries and rollback if any of the queries fail or commit if all queries are
// successful.
func (s *User) CreateManagementSession() error {
	return s.client.createManagementSession()
}

// EndManagementSession ends the session created to execute a series of actions
// taken to manage an user. Ending a session makes it possible to
// rollback if any of the queries in a management session fail or commit if all
// queries are successful.
func (s *User) EndManagementSession() error {
	return s.client.endManagementSession()
}

// Create a new user.
func (s *User) Create(u *user.User) (e error) {
	query, e := NewQuery(u)
	if e != nil {
		return e
	}

	userInsertQuery := query.Create(schema, userTable)
	userInsertQuery = s.client.db.Rebind(userInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(u, userInsertQuery, queryParams...)
	} else {
		e = s.client.db.Get(u, userInsertQuery, queryParams...)
	}
	var uniqueViolation pq.ErrorCode = "23505"
	if pqError, ok := e.(*pq.Error); e != nil && !ok {
		return e
	} else if pqError != nil && pqError.Code == uniqueViolation {
		return user.ErrUserExists
	} else if pqError != nil {
		return pqError
	}

	return nil
}

// Read a user by uuid.
func (s *User) Read(u *user.User, byEmail bool) (res *user.User, e error) {
	filter := "UUID=?"
	queryParam := u.UUID
	if byEmail {
		filter = "EMAIL=?"
		queryParam = u.Email
	}

	query, err := NewQuery(u)
	if err != nil {
		return nil, err
	}

	userSelectQuery := query.Read(schema, userTable, filter)
	userSelectQuery = s.client.db.Rebind(userSelectQuery)

	res = &user.User{}
	if err := s.client.db.Get(res, userSelectQuery, queryParam); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil

}

// Update the , searching by email or uuid.
func (s *User) Update(u *user.User, byEmail bool) (e error) {
	filter := "UUID=?"
	queryParam := u.UUID
	if byEmail {
		filter = "EMAIL=?"
		queryParam = u.Email
	}

	query, e := NewQuery(u)
	if e != nil {
		return e
	}

	userUpdateQuery := query.Update(schema, userTable, filter)
	userUpdateQuery = s.client.db.Rebind(userUpdateQuery)

	includeNil := true
	queryParams := append(query.ModelValues(includeNil), queryParam)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(u, userUpdateQuery, queryParams...)
	} else {
		e = s.client.db.Get(u, userUpdateQuery, queryParams...)
	}
	if e != nil {
		return e
	}

	return nil
}

// Delete a user searching by uuid or email.
func (s *User) Delete(u *user.User, byEmail bool) (e error) {
	filter := "UUID=?"
	queryParam := u.UUID
	if byEmail {
		filter = "EMAIL=?"
		queryParam = u.Email
	}

	query, e := NewQuery(u)
	if e != nil {
		return e
	}
	userDeleteQuery := query.Delete(schema, userTable, filter)
	userDeleteQuery = s.client.db.Rebind(userDeleteQuery)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(u, userDeleteQuery, queryParam)
	} else {
		e = s.client.db.Get(u, userDeleteQuery, queryParam)
	}
	if e != nil && e == sql.ErrNoRows {
		return user.ErrUserNotFound
	} else if e != nil {
		return e
	}

	return nil
}
