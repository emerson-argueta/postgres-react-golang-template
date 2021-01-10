package persistence

import (
	"context"
	"database/sql"

	"emersonargueta/m/v1/modules/identity/domain/user"
	"emersonargueta/m/v1/modules/identity/repository"
	"emersonargueta/m/v1/shared/infrastructure/database"
)

var _ repository.UserRepo = &User{}

// User represents a service for managing a user.
type User struct {
	Client database.Client
}

// CreateUser if successful. If the user
// exists, returns ErrUserExists.
func (s *User) CreateUser(u user.User) (e error) {
	userPersistenceModel := UserDomainToPersistence(u)

	userInsertQuery := s.Client.Query().Create(userPersistenceModel, IdentitySchema, UserTable)
	userInsertQuery = s.Client.DB().Rebind(userInsertQuery)

	includeNil := true
	queryParams := s.Client.Query().ModelValues(userPersistenceModel, includeNil)

	ctx := context.WithValue(context.Background(), "aggregateid", u.GetID())
	_, e = s.Client.DB().ExecContext(ctx, userInsertQuery, queryParams...)

	return e

}

// RetrieveUserByEmail searching by email. If the user does not exists,
// returns ErrUserNotFound.
func (s *User) RetrieveUserByEmail(email user.Email) (res user.User, e error) {
	filter := "EMAIL=?"
	queryParam := email.ToString()

	userSelectQuery := s.Client.Query().Read(IdentitySchema, UserTable, filter)
	userSelectQuery = s.Client.DB().Rebind(userSelectQuery)

	dbRes := &UserDTO{}
	e = s.Client.DB().Get(dbRes, userSelectQuery, queryParam)
	if e == sql.ErrNoRows || dbRes == nil {
		return nil, user.ErrUserNotFound
	}

	return UserPersistenceToDomain(*dbRes)

}

// RetrieveUserByID searching by id. If the user does not exists,
// returns ErrUserNotFound.
func (s *User) RetrieveUserByID(id string) (res user.User, e error) {
	filter := "ID=?"
	queryParam := id

	userSelectQuery := s.Client.Query().Read(IdentitySchema, UserTable, filter)
	userSelectQuery = s.Client.DB().Rebind(userSelectQuery)

	dbRes := &UserDTO{}
	e = s.Client.DB().Get(dbRes, userSelectQuery, queryParam)
	if e == sql.ErrNoRows || dbRes == nil {
		return nil, user.ErrUserNotFound
	}

	return UserPersistenceToDomain(*dbRes)

}

// UpdateUser searching by uuid. If the user does not exists, returns
// ErrUserNotFound.
func (s *User) UpdateUser(u user.User) (e error) {
	filter := "ID=?"

	userPersistenceModel := UserDomainToPersistence(u)
	userUpdateQuery := s.Client.Query().Update(userPersistenceModel, IdentitySchema, UserTable, filter)
	userUpdateQuery = s.Client.DB().Rebind(userUpdateQuery)

	includeNil := true
	queryParams := append(s.Client.Query().ModelValues(userPersistenceModel, includeNil), u.GetID())

	e = s.Client.DB().Get(userPersistenceModel, userUpdateQuery, queryParams...)
	if e == sql.ErrNoRows {
		return user.ErrUserNotFound
	}

	return e
}

// DeleteUser searching by uuid. If the user does not exists, returns
// ErrUserNotFound.
func (s *User) DeleteUser(uuid string) (e error) {
	filter := "ID=?"

	userDeleteQuery := s.Client.Query().Delete(IdentitySchema, UserTable, filter)
	userDeleteQuery = s.Client.DB().Rebind(userDeleteQuery)

	userPersistenceModel := &UserDTO{}
	queryParam := uuid
	e = s.Client.DB().Get(userPersistenceModel, userDeleteQuery, queryParam)
	if e == sql.ErrNoRows {
		return user.ErrUserNotFound
	}

	return e
}
