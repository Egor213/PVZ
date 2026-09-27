package order

import (
	"context"
	"fmt"
)

type CancelOrderInput struct {
	ID string
}

type CancelOrderOutput struct {
	ID        string
	Status    Status
	UpdatedAt int64
}

type CancelOrder struct {
	repo     Repository
	notifier Notifier
}

func NewCancelOrder(repo Repository, notifier Notifier) *CancelOrder {
	return &CancelOrder{
		repo:     repo,
		notifier: notifier,
	}
}

func (uc *CancelOrder) Execute(ctx context.Context, in CancelOrderInput) (CancelOrderOutput, error) {
	o, err := uc.repo.Get(ctx, in.ID)
	if err != nil {
		return CancelOrderOutput{}, fmt.Errorf("get order: %w", err)
	}

	if err := o.Cancel(); err != nil {
		return CancelOrderOutput{}, err
	}

	if err := uc.repo.Save(ctx, o); err != nil {
		return CancelOrderOutput{}, fmt.Errorf("save order: %w", err)
	}

	if uc.notifier != nil {
		if err := uc.notifier.OrderCancelled(ctx, o); err != nil {
			// log error but don't fail
		}
	}

	return CancelOrderOutput{
		ID:        o.ID,
		Status:    o.Status,
		UpdatedAt: o.UpdatedAt,
	}, nil
}