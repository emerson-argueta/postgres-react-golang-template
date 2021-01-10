package domain

import "fmt"

var _ AggregateRoot = &AbstractAggregateRoot{}

type AggregateRoot interface {
	GetID() string
	GetDomainEvents() []Event
	AddDomainEvent(domainEvent Event)
	ClearEvents()
	LogDomainEventAdded(domainEvent Event)
}
type AbstractAggregateRoot struct {
	DomainEvents []Event
	Name         string
	ID           string
}

func (a *AbstractAggregateRoot) GetID() string {
	return a.ID
}

func (a *AbstractAggregateRoot) GetDomainEvents() []Event {
	return a.DomainEvents
}

func (a *AbstractAggregateRoot) AddDomainEvent(domainEvent Event) {

	a.DomainEvents = append(a.DomainEvents, domainEvent)
	// Add this aggregate instance to the domain event's list of aggregates who's
	// events it eventually needs to dispatch.
	DomainEvents.MarkAggregateForDispatch(a)
	// Log the domain event
	a.LogDomainEventAdded(domainEvent)
}

func (a *AbstractAggregateRoot) ClearEvents() {
	a.DomainEvents = make([]Event, 0)
}

func (a *AbstractAggregateRoot) LogDomainEventAdded(domainEvent Event) {
	aggregateName := a.Name
	eventName := domainEvent.GetName()
	fmt.Printf("[Domain Event Created]: %s ==> %s \n", aggregateName, eventName)
}
