package order

import "context"

type IDGenerator interface {
	New() string
}

type Clock interface {
	Now() int64
}

type Notifier interface {
	OrderCreated(ctx context.Context, order *Order) error
	OrderCancelled(ctx context.Context, order *Order) error
}