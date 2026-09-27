package payment

import "context"

type IDGenerator interface {
	New() string
}

type Clock interface {
	Now() int64
}

type Notifier interface {
	PaymentCreated(ctx context.Context, payment *Payment) error
	PaymentUpdated(ctx context.Context, payment *Payment) error
}