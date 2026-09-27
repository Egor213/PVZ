package order

import "context"

type Repository interface {
	Save(ctx context.Context, o *Order) error
	Get(ctx context.Context, id string) (*Order, error)
	GetAll(ctx context.Context) ([]*Order, error)
	Delete(ctx context.Context, id string) error
}