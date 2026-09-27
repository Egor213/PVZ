package payment

import (
	"context"
	"fmt"
	"time"
)

type UpdatePaymentInput struct {
	ID     string
	Status *Status
}

type UpdatePaymentOutput struct {
	ID        string
	OrderID   string
	Amount    int64
	Status    Status
	CreatedAt int64
	UpdatedAt int64
}

type UpdatePayment struct {
	repo Repository
}

func NewUpdatePayment(repo Repository) *UpdatePayment {
	return &UpdatePayment{repo: repo}
}

func (uc *UpdatePayment) Execute(ctx context.Context, in UpdatePaymentInput) (UpdatePaymentOutput, error) {
	p, err := uc.repo.Get(ctx, in.ID)
	if err != nil {
		return UpdatePaymentOutput{}, fmt.Errorf("get payment: %w", err)
	}

	if in.Status != nil {
		if !in.Status.Valid() {
			return UpdatePaymentOutput{}, ErrInvalidStatus
		}
		p.Status = *in.Status
		p.UpdatedAt = time.Now().Unix()
	}

	if err := uc.repo.Save(ctx, p); err != nil {
		return UpdatePaymentOutput{}, fmt.Errorf("save payment: %w", err)
	}

	return UpdatePaymentOutput{
		ID:        p.ID,
		OrderID:   p.OrderID,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}