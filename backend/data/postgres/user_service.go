package postgres

import (
	"database/sql"

	"emersonargueta/m/v1/user"

	"github.com/lib/pq"
)

const userTable = "user"
const schema = "identity"

var _ user.Service = &User{}

// User represents a service for managing a user.
type User struct {
	client *Client
}

// Register a new user.
func (s *User) Register(u *user.User) (e error) {
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

// Retrieve a user by uuid.
func (s *User) Retrieve(u *user.User, byEmail bool) (res *user.User, e error) {
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

// UnRegister a user searching by uuid or email.
func (s *User) UnRegister(u *user.User, byEmail bool) (e error) {
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

// LookUpDomain by domain's name
func (s *User) LookUpDomain(domain *user.Domain) (res *user.Domain, e error) {
	filter := "NAME=?"
	queryParam := domain.Name

	query, err := NewQuery(domain)
	if err != nil {
		return nil, err
	}

	userSelectQuery := query.Read(schema, userTable, filter)
	userSelectQuery = s.client.db.Rebind(userSelectQuery)

	res = &user.Domain{}
	if err := s.client.db.Get(res, userSelectQuery, queryParam); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}
