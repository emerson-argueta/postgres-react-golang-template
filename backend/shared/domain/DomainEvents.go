package domain

// DomainEvents for all aggregate roots
var DomainEvents = &domainEvents{
	handlersMap:      make(map[string][]func(Event)),
	markedAggregates: make([]AggregateRoot, 0),
}

// Abilities of domain events
type Abilities interface {
	MarkAggregateForDispatch(aggregate AggregateRoot)
	DispatchEventsForAggregate(id string)
	Register(callback func(Event), eventName string)
	ClearHandlers()
	ClearMarkedAggregates()
}

type domainEvents struct {
	Abilities
	handlersMap      map[string][]func(Event)
	markedAggregates []AggregateRoot
}

// MarkAggregateForDispatch called by aggregate root objects that have created
// domain events to eventually be dispatched when the infrastructure commits the
// unit of work.
func (a *domainEvents) MarkAggregateForDispatch(aggregate AggregateRoot) {
	id := aggregate.GetID()
	aggregateFound := a.findMarkedAggregateByID(id)
	if aggregateFound == nil {
		a.markedAggregates = append(a.markedAggregates, aggregate)
	}
}
func (a *domainEvents) findMarkedAggregateByID(id string) AggregateRoot {
	var found AggregateRoot = nil
	for _, aggregate := range a.markedAggregates {
		if aggregate.GetID() == id {
			found = aggregate
		}
	}

	return found
}

// DispatchEventsForAggregate from a list of marked aggregate events. Will
// trigger a dispatch to notify that an domain event has occured.
func (a *domainEvents) DispatchEventsForAggregate(id string) {
	aggregate := a.findMarkedAggregateByID(id)

	if aggregate != nil {
		a.dispatchAggregateEvents(aggregate)
		aggregate.ClearEvents()
		a.removeAggregateFromMarkedDispatchList(aggregate)
	}
}
func (a *domainEvents) dispatchAggregateEvents(aggregate AggregateRoot) {
	for _, domainEvent := range aggregate.GetDomainEvents() {
		a.dispatch(domainEvent)
	}
}
func (a *domainEvents) dispatch(event Event) {
	eventName := event.GetName()

	if handlers, exists := a.handlersMap[eventName]; exists {
		for _, handler := range handlers {
			handler(event)
		}
	}
}
func (a *domainEvents) removeAggregateFromMarkedDispatchList(aggregate AggregateRoot) {
	var index int
	for i := range a.markedAggregates {
		if a.markedAggregates[i].GetID() == aggregate.GetID() {
			index = i
		}
	}

	removeAtIndex(a.markedAggregates, index)
}
func removeAtIndex(arr []AggregateRoot, index int) {
	arr[index] = arr[len(arr)-1]
	arr = arr[:len(arr)-1]
}

// Register to a domain event that will trigger a callback function
func (a *domainEvents) Register(callback func(event Event), eventName string) {
	if _, exists := a.handlersMap[eventName]; !exists {
		a.handlersMap[eventName] = make([]func(Event), 0)
	}
	a.handlersMap[eventName] = append(a.handlersMap[eventName], callback)
}

// ClearHandlers when they are no longer necessary.
func (a *domainEvents) ClearHandlers() {
	a.handlersMap = make(map[string][]func(Event))
}

// ClearMarkedAggregates when they are no longer necessary.
func (a *domainEvents) ClearMarkedAggregates() {
	a.markedAggregates = nil
}
