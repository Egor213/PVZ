package payment

import (
	"context"
	"fmt"
)

type RefundPaymentInput struct {
	ID string
}

type RefundPaymentOutput struct {
	ID        string
	OrderID   string
	Amount    int64
	Status    Status
	CreatedAt int64
	UpdatedAt int64
}

type RefundPayment struct {
	repo     Repository
	notifier Notifier
}

func NewRefundPayment(repo Repository, notifier Notifier) *RefundPayment {
	return &RefundPayment{
		repo:     repo,
		notifier: notifier,
	}
}

func (uc *RefundPayment) Execute(ctx context.Context, in RefundPaymentInput) (RefundPaymentOutput, error) {
	p, err := uc.repo.Get(ctx, in.ID)
	if err != nil {
		return RefundPaymentOutput{}, fmt.Errorf("get payment: %w", err)
	}

	if err := p.Refund(); err != nil {
		return RefundPaymentOutput{}, err
	}

	if err := uc.repo.Save(ctx, p); err != nil {
		return RefundPaymentOutput{}, fmt.Errorf("save payment: %w", err)
	}

	if uc.notifier != nil {
		if err := uc.notifier.PaymentUpdated(ctx, p); err != nil {
			// log error but don't fail
		}
	}

	return RefundPaymentOutput{
		ID:        p.ID,
		OrderID:   p.OrderID,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}