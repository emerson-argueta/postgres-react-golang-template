package postgres

import (
	"database/sql"
	"emersonargueta/m/v1/identity"
	"emersonargueta/m/v1/identity/user"

	"github.com/lib/pq"
)

var _ user.Processes = &userservice{}

// userservice represents a service for managing a user.
type userservice struct {
	client *Client
}

// CreateUser if successful. If the user
// exists, returns ErrUserExists.
func (s *userservice) CreateUser(u *user.User) (res *user.User, e error) {
	query, e := NewQuery(u)
	if e != nil {
		return nil, e
	}

	userInsertQuery := query.Create(IdentitySchema, UserTable)
	userInsertQuery = s.client.db.Rebind(userInsertQuery)

	includeNil := true
	queryParams := query.ModelValues(includeNil)

	res = &user.User{}

	e = s.client.db.Get(res, userInsertQuery, queryParams...)

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

// RetrieveUser searching by email. If the user does not exists,
// returns ErrUserNotFound.
func (s *userservice) RetrieveUser(email string) (res *user.User, e error) {
	filter := "EMAIL=?"
	queryParam := email

	query, e := NewQuery(&user.User{})
	if e != nil {
		return nil, e
	}

	userSelectQuery := query.Read(IdentitySchema, UserTable, filter)
	userSelectQuery = s.client.db.Rebind(userSelectQuery)

	res = &user.User{}
	e = s.client.db.Get(res, userSelectQuery, queryParam)
	if e == sql.ErrNoRows {
		return nil, identity.ErrUserNotFound
	}

	return res, e

}

// UpdateUser searching by uuid. If the user does not exists, returns
// ErrUserNotFound.
func (s *userservice) UpdateUser(u *user.User) (e error) {
	filter := "UUID=?"
	queryParam := u.UUID

	query, e := NewQuery(u)
	if e != nil {
		return e
	}

	userUpdateQuery := query.Update(IdentitySchema, UserTable, filter)
	userUpdateQuery = s.client.db.Rebind(userUpdateQuery)

	includeNil := true
	queryParams := append(query.ModelValues(includeNil), queryParam)

	e = s.client.db.Get(u, userUpdateQuery, queryParams...)
	if e == sql.ErrNoRows {
		return identity.ErrUserNotFound
	}

	return e
}

// DeleteUser searching by uuid. If the user does not exists, returns
// ErrUserNotFound.
func (s *userservice) DeleteUser(uuid string) (e error) {
	filter := "UUID=?"
	queryParam := uuid

	query, e := NewQuery(&user.User{})
	if e != nil {
		return e
	}
	userDeleteQuery := query.Delete(IdentitySchema, UserTable, filter)
	userDeleteQuery = s.client.db.Rebind(userDeleteQuery)

	e = s.client.db.Get(&user.User{}, userDeleteQuery, queryParam)
	if e == sql.ErrNoRows {
		return identity.ErrUserNotFound
	}

	return e
}
