package postgres

import (
	"database/sql"
	"emersonargueta/m/v1/identity"
	"emersonargueta/m/v1/identity/user"

	"github.com/lib/pq"
)

// UserTable stores user information for the identity domain
const UserTable = "user"

// Schema used to group tables used in the identity domain
const Schema = "identity"

var _ user.Service = &User{}

// User represents a service for managing a user.
type User struct {
	client *Client
}

// CreateUser a new user.
func (s *User) CreateUser(u *user.User) (e error) {
	query, e := NewQuery(u)
	if e != nil {
		return e
	}

	userInsertQuery := query.Create(Schema, UserTable)
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
		return identity.ErrUserExists
	} else if pqError != nil {
		return pqError
	}

	return nil
}

// RetrieveUser a user by uuid or email.
func (s *User) RetrieveUser(u *user.User, byEmail bool) (res *user.User, e error) {
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

	userSelectQuery := query.Read(Schema, UserTable, filter)
	userSelectQuery = s.client.db.Rebind(userSelectQuery)

	res = &user.User{}
	if err := s.client.db.Get(res, userSelectQuery, queryParam); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil

}

// UpdateUser searching by email or uuid.
func (s *User) UpdateUser(u *user.User, byEmail bool) (e error) {
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

	userUpdateQuery := query.Update(Schema, UserTable, filter)
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

// DeleteUser searching by uuid or email.
func (s *User) DeleteUser(u *user.User, byEmail bool) (e error) {
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
	userDeleteQuery := query.Delete(Schema, UserTable, filter)
	userDeleteQuery = s.client.db.Rebind(userDeleteQuery)

	if s.client.transaction != nil {
		e = s.client.transaction.Get(u, userDeleteQuery, queryParam)
	} else {
		e = s.client.db.Get(u, userDeleteQuery, queryParam)
	}
	if e != nil && e == sql.ErrNoRows {
		return identity.ErrUserNotFound
	} else if e != nil {
		return e
	}

	return nil
}
