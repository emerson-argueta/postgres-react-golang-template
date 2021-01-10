package achiever

import (
	"emersonargueta/m/v1/shared/domain"
	"time"
)

// AchieverCreatedEventName of domain event
const AchieverCreatedEventName = "AchieverCreated"

type AchieverCreated interface {
	GetAchiever() Achiever
	domain.Event
}
type achieverCreated struct {
	Achiever Achiever
	event    *domain.AbstractEvent
}

func NewAchieverCreated(achiever Achiever) AchieverCreated {
	achieverCreated := &achieverCreated{
		Achiever: achiever,
		event:    &domain.AbstractEvent{},
	}

	achieverCreated.event.TimeOccured = time.Now()
	achieverCreated.event.AggregateID = achiever.GetID()
	achieverCreated.event.Name = AchieverCreatedEventName

	return achieverCreated
}
func (e *achieverCreated) GetAggregateID() string {
	return e.event.AggregateID
}
func (e *achieverCreated) GetName() string {
	return e.event.Name
}
func (e *achieverCreated) GetAchiever() Achiever {
	return e.Achiever
}
