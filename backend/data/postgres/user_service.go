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
func (s *User) CreateUser(u *user.User) (res *user.User, e error) {
	query, e := NewQuery(u)
	if e != nil {
		return nil, e
	}

	userInsertQuery := query.Create(Schema, UserTable)
	userInsertQuery = s.client.db.Rebind(userInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	res = &user.User{}

	e = s.client.db.Get(u, userInsertQuery, queryParams...)

	var uniqueViolation pq.ErrorCode = "23505"
	if pqError, ok := e.(*pq.Error); e != nil && !ok {
		return nil, e
	} else if pqError != nil && pqError.Code == uniqueViolation {
		return nil, identity.ErrUserExists
	} else if pqError != nil {
		return nil, pqError
	}

	return res, e
}

// RetrieveUser a user by email.
func (s *User) RetrieveUser(email string) (res *user.User, e error) {
	filter := "EMAIL=?"
	queryParam := email

	query, e := NewQuery(&user.User{})
	if e != nil {
		return nil, e
	}

	userSelectQuery := query.Read(Schema, UserTable, filter)
	userSelectQuery = s.client.db.Rebind(userSelectQuery)

	res = &user.User{}
	e = s.client.db.Get(res, userSelectQuery, queryParam)
	if e == sql.ErrNoRows {
		return nil, identity.ErrUserNotFound
	}

	return res, e

}

// UpdateUser searching by uuid.
func (s *User) UpdateUser(u *user.User) (e error) {
	filter := "UUID=?"
	queryParam := u.UUID

	query, e := NewQuery(u)
	if e != nil {
		return e
	}

	userUpdateQuery := query.Update(Schema, UserTable, filter)
	userUpdateQuery = s.client.db.Rebind(userUpdateQuery)

	includeNil := true
	queryParams := append(query.ModelValues(includeNil), queryParam)

	e = s.client.db.Get(u, userUpdateQuery, queryParams...)
	if e == sql.ErrNoRows {
		return identity.ErrUserNotFound
	}

	return e
}

// DeleteUser searching by uuid.
func (s *User) DeleteUser(uuid string) (e error) {
	filter := "UUID=?"
	queryParam := uuid

	query, e := NewQuery(&user.User{})
	if e != nil {
		return e
	}
	userDeleteQuery := query.Delete(Schema, UserTable, filter)
	userDeleteQuery = s.client.db.Rebind(userDeleteQuery)

	e = s.client.db.Get(&user.User{}, userDeleteQuery, queryParam)
	if e == sql.ErrNoRows {
		return identity.ErrUserNotFound
	}

	return e
}
