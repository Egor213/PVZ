package order

import (
	"context"
	"fmt"
)

type UpdateOrderInput struct {
	ID          string
	Title       *string
	Description *string
	Status      *Status
}

type UpdateOrderOutput struct {
	ID          string
	Title       string
	Description string
	Status      Status
	CreatedAt   int64
	UpdatedAt   int64
}

type UpdateOrder struct {
	repo Repository
}

func NewUpdateOrder(repo Repository) *UpdateOrder {
	return &UpdateOrder{repo: repo}
}

func (uc *UpdateOrder) Execute(ctx context.Context, in UpdateOrderInput) (UpdateOrderOutput, error) {
	o, err := uc.repo.Get(ctx, in.ID)
	if err != nil {
		return UpdateOrderOutput{}, fmt.Errorf("get order: %w", err)
	}

	if in.Title != nil {
		if err := o.UpdateTitle(*in.Title); err != nil {
			return UpdateOrderOutput{}, err
		}
	}
	if in.Description != nil {
		o.UpdateDescription(*in.Description)
	}
	if in.Status != nil {
		if err := o.UpdateStatus(*in.Status); err != nil {
			return UpdateOrderOutput{}, err
		}
	}

	if err := uc.repo.Save(ctx, o); err != nil {
		return UpdateOrderOutput{}, fmt.Errorf("save order: %w", err)
	}

	return UpdateOrderOutput{
		ID:          o.ID,
		Title:       o.Title,
		Description: o.Description,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}, nil
}