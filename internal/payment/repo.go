package payment

import "context"

type Repository interface {
	Save(ctx context.Context, p *Payment) error
	Get(ctx context.Context, id string) (*Payment, error)
	GetByOrderID(ctx context.Context, orderID string) ([]*Payment, error)
	GetAll(ctx context.Context) ([]*Payment, error)
	Delete(ctx context.Context, id string) error
}