package user

import (
	"emersonargueta/m/v1/shared/domain"
	"time"
)

// UserCreatedEventName of domain event
const UserCreatedEventName = "UserCreated"

type UserCreated interface {
	GetUser() User
	domain.Event
}
type userCreated struct {
	User  User
	event *domain.AbstractEvent
}

func NewUserCreated(user User) UserCreated {
	userCreated := &userCreated{
		User:  user,
		event: &domain.AbstractEvent{},
	}

	userCreated.event.TimeOccured = time.Now()
	userCreated.event.AggregateID = user.GetID()
	userCreated.event.Name = UserCreatedEventName

	return userCreated
}
func (e *userCreated) GetAggregateID() string {
	return e.event.AggregateID
}
func (e *userCreated) GetName() string {
	return e.event.Name
}
func (e *userCreated) GetUser() User {
	return e.User
}
