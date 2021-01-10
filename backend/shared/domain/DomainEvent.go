package domain

import "time"

// Event for domain aggregate
type Event interface {
	GetAggregateID() string
	GetName() string
}

// AbstractEvent for domains created by aggregate roots
type AbstractEvent struct {
	TimeOccured time.Time
	AggregateID string
	Name        string
}
