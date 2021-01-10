package user

import (
	"emersonargueta/m/v1/shared/domain"

	"github.com/google/uuid"
)

// Fields for user
type Fields struct {
	Role         *Role
	Email        *Email
	HashPassword *HashPassword
}

// User is part of the identity subdomain to support core domains which need it
type User interface {
	GetEmail() Email
	GetHashPassword() HashPassword
	GetRole() Role
	GetID() string
	SetEmail(Email)
	SetHashPassword(HashPassword)
	SetRole(Role)
}

//User model
type user struct {
	ID           string
	Role         Role
	Email        Email
	HashPassword HashPassword
	aggregate    *domain.AbstractAggregateRoot
}

// Create user with role, email and password
func Create(userFields *Fields, id *string) (res User, e error) {
	if userFields.Email == nil || userFields.HashPassword == nil || userFields.Role == nil {
		return nil, ErrUserIncompleteDetails
	}

	user := &user{
		Role:         *userFields.Role,
		Email:        *userFields.Email,
		HashPassword: *userFields.HashPassword,
	}

	isNewUser := (id == nil)
	if isNewUser {
		user.ID = uuid.New().String()
	}
	if !isNewUser {
		user.ID = *id
	}

	user.aggregate = &domain.AbstractAggregateRoot{}
	user.aggregate.DomainEvents = make([]domain.Event, 0)
	user.aggregate.Name = "User"
	user.aggregate.ID = user.ID

	if isNewUser {
		user.aggregate.AddDomainEvent(NewUserCreated(user))
	}

	return user, nil
}

func (u *user) GetEmail() Email {
	return u.Email
}
func (u *user) GetHashPassword() HashPassword {
	return u.HashPassword
}
func (u *user) GetRole() Role {
	return u.Role
}
func (u *user) GetID() string {
	return u.ID
}
func (u *user) SetEmail(email Email) {
	u.Email = email
}
func (u *user) SetHashPassword(hashPassword HashPassword) {
	u.HashPassword = hashPassword
}
func (u *user) SetRole(role Role) {
	u.Role = role
}
