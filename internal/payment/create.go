package payment

import (
	"context"
	"fmt"
)

type CreatePaymentInput struct {
	OrderID string
	Amount  int64
}

type CreatePaymentOutput struct {
	ID        string
	OrderID   string
	Amount    int64
	Status    Status
	CreatedAt int64
	UpdatedAt int64
}

type CreatePayment struct {
	repo     Repository
	ids      IDGenerator
	clock    Clock
	notifier Notifier
}

func NewCreatePayment(repo Repository, ids IDGenerator, clock Clock, notifier Notifier) *CreatePayment {
	return &CreatePayment{
		repo:     repo,
		ids:      ids,
		clock:    clock,
		notifier: notifier,
	}
}

func (uc *CreatePayment) Execute(ctx context.Context, in CreatePaymentInput) (CreatePaymentOutput, error) {
	p, err := NewPayment(in.OrderID, in.Amount)
	if err != nil {
		return CreatePaymentOutput{}, err
	}
	p.ID = uc.ids.New()
	p.CreatedAt = uc.clock.Now()
	p.UpdatedAt = uc.clock.Now()

	if err := uc.repo.Save(ctx, p); err != nil {
		return CreatePaymentOutput{}, fmt.Errorf("save payment: %w", err)
	}

	if uc.notifier != nil {
		if err := uc.notifier.PaymentCreated(ctx, p); err != nil {
			// log error but don't fail
		}
	}

	return CreatePaymentOutput{
		ID:        p.ID,
		OrderID:   p.OrderID,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}