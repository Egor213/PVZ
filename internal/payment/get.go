package payment

import (
	"context"
	"fmt"
)

type GetPaymentInput struct {
	ID string
}

type GetPaymentOutput struct {
	ID        string
	OrderID   string
	Amount    int64
	Status    Status
	CreatedAt int64
	UpdatedAt int64
}

type GetPayment struct {
	repo Repository
}

func NewGetPayment(repo Repository) *GetPayment {
	return &GetPayment{repo: repo}
}

func (uc *GetPayment) Execute(ctx context.Context, in GetPaymentInput) (GetPaymentOutput, error) {
	p, err := uc.repo.Get(ctx, in.ID)
	if err != nil {
		return GetPaymentOutput{}, fmt.Errorf("get payment: %w", err)
	}
	return GetPaymentOutput{
		ID:        p.ID,
		OrderID:   p.OrderID,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}

type GetPaymentsByOrderInput struct {
	OrderID string
}

type GetPaymentsByOrderOutput struct {
	Payments []GetPaymentOutput
}

type GetPaymentsByOrder struct {
	repo Repository
}

func NewGetPaymentsByOrder(repo Repository) *GetPaymentsByOrder {
	return &GetPaymentsByOrder{repo: repo}
}

func (uc *GetPaymentsByOrder) Execute(ctx context.Context, in GetPaymentsByOrderInput) (GetPaymentsByOrderOutput, error) {
	payments, err := uc.repo.GetByOrderID(ctx, in.OrderID)
	if err != nil {
		return GetPaymentsByOrderOutput{}, fmt.Errorf("get payments by order: %w", err)
	}
	out := make([]GetPaymentOutput, len(payments))
	for i, p := range payments {
		out[i] = GetPaymentOutput{
			ID:        p.ID,
			OrderID:   p.OrderID,
			Amount:    p.Amount,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}
	return GetPaymentsByOrderOutput{Payments: out}, nil
}

type GetAllPaymentsInput struct{}

type GetAllPaymentsOutput struct {
	Payments []GetPaymentOutput
}

type GetAllPayments struct {
	repo Repository
}

func NewGetAllPayments(repo Repository) *GetAllPayments {
	return &GetAllPayments{repo: repo}
}

func (uc *GetAllPayments) Execute(ctx context.Context, in GetAllPaymentsInput) (GetAllPaymentsOutput, error) {
	payments, err := uc.repo.GetAll(ctx)
	if err != nil {
		return GetAllPaymentsOutput{}, fmt.Errorf("get all payments: %w", err)
	}
	out := make([]GetPaymentOutput, len(payments))
	for i, p := range payments {
		out[i] = GetPaymentOutput{
			ID:        p.ID,
			OrderID:   p.OrderID,
			Amount:    p.Amount,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}
	return GetAllPaymentsOutput{Payments: out}, nil
}