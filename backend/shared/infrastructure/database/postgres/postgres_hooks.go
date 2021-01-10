package postgres

import (
	"context"
	"emersonargueta/m/v1/shared/domain"
	"fmt"

	"github.com/gchaincl/sqlhooks"
)

var _ sqlhooks.Hooks = &Hooks{}

// Hooks satisfies the sqlhook.Hooks interface
type Hooks struct{}

// Before hook satisfies hooks contract
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	fmt.Printf("> %s %q", query, args)

	return ctx, nil
}

// After hook will dipatch domain events if context has aggregateid key
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	ctxID := ctx.Value("aggregateid")
	if ctxID != nil {
		aggregateID := ctxID.(string)
		domain.DomainEvents.DispatchEventsForAggregate(aggregateID)
	}
	return ctx, nil
}
