package order

import (
	"context"
	"fmt"
)

type CreateOrderInput struct {
	Title       string
	Description string
}

type CreateOrderOutput struct {
	ID          string
	Title       string
	Description string
	Status      Status
	CreatedAt   int64
	UpdatedAt   int64
}

type CreateOrder struct {
	repo     Repository
	ids      IDGenerator
	clock    Clock
	notifier Notifier
}

func NewCreateOrder(repo Repository, ids IDGenerator, clock Clock, notifier Notifier) *CreateOrder {
	return &CreateOrder{
		repo:     repo,
		ids:      ids,
		clock:    clock,
		notifier: notifier,
	}
}

func (uc *CreateOrder) Execute(ctx context.Context, in CreateOrderInput) (CreateOrderOutput, error) {
	o, err := NewOrder(in.Title, in.Description)
	if err != nil {
		return CreateOrderOutput{}, err
	}
	o.ID = uc.ids.New()
	o.CreatedAt = uc.clock.Now()
	o.UpdatedAt = uc.clock.Now()

	if err := uc.repo.Save(ctx, o); err != nil {
		return CreateOrderOutput{}, fmt.Errorf("save order: %w", err)
	}

	if uc.notifier != nil {
		if err := uc.notifier.OrderCreated(ctx, o); err != nil {
			// log error but don't fail the order creation
		}
	}

	return CreateOrderOutput{
		ID:          o.ID,
		Title:       o.Title,
		Description: o.Description,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}, nil
}