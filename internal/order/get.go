package order

import (
	"context"
	"fmt"
)

type GetOrderInput struct {
	ID string
}

type GetOrderOutput struct {
	ID          string
	Title       string
	Description string
	Status      Status
	CreatedAt   int64
	UpdatedAt   int64
}

type GetOrder struct {
	repo Repository
}

func NewGetOrder(repo Repository) *GetOrder {
	return &GetOrder{repo: repo}
}

func (uc *GetOrder) Execute(ctx context.Context, in GetOrderInput) (GetOrderOutput, error) {
	o, err := uc.repo.Get(ctx, in.ID)
	if err != nil {
		return GetOrderOutput{}, fmt.Errorf("get order: %w", err)
	}
	return GetOrderOutput{
		ID:          o.ID,
		Title:       o.Title,
		Description: o.Description,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}, nil
}

type GetAllOrdersInput struct{}

type GetAllOrdersOutput struct {
	Orders []GetOrderOutput
}

type GetAllOrders struct {
	repo Repository
}

func NewGetAllOrders(repo Repository) *GetAllOrders {
	return &GetAllOrders{repo: repo}
}

func (uc *GetAllOrders) Execute(ctx context.Context, in GetAllOrdersInput) (GetAllOrdersOutput, error) {
	orders, err := uc.repo.GetAll(ctx)
	if err != nil {
		return GetAllOrdersOutput{}, fmt.Errorf("get all orders: %w", err)
	}
	out := make([]GetOrderOutput, len(orders))
	for i, o := range orders {
		out[i] = GetOrderOutput{
			ID:          o.ID,
			Title:       o.Title,
			Description: o.Description,
			Status:      o.Status,
			CreatedAt:   o.CreatedAt,
			UpdatedAt:   o.UpdatedAt,
		}
	}
	return GetAllOrdersOutput{Orders: out}, nil
}