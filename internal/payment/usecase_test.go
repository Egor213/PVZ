package payment

import (
	"context"
	"errors"
	"testing"
)

type fakeRepo struct {
	saved *Payment
}

func (r *fakeRepo) Save(ctx context.Context, p *Payment) error {
	r.saved = p
	return nil
}

func (r *fakeRepo) Get(ctx context.Context, id string) (*Payment, error) {
	if r.saved != nil && r.saved.ID == id {
		return r.saved, nil
	}
	return nil, ErrNotFound
}

func (r *fakeRepo) GetByOrderID(ctx context.Context, orderID string) ([]*Payment, error) {
	return nil, nil
}

func (r *fakeRepo) GetAll(ctx context.Context) ([]*Payment, error) {
	return nil, nil
}

func (r *fakeRepo) Delete(ctx context.Context, id string) error {
	return nil
}

type fixedIDGen struct{ id string }

func (g *fixedIDGen) New() string { return g.id }

type fixedClock struct{ now int64 }

func (c *fixedClock) Now() int64 { return c.now }

func TestCreatePayment_Success(t *testing.T) {
	repo := &fakeRepo{}
	ids := &fixedIDGen{id: "payment-1"}
	clock := &fixedClock{now: 2000}
	uc := NewCreatePayment(repo, ids, clock, nil)

	out, err := uc.Execute(context.Background(), CreatePaymentInput{
		OrderID: "order-1",
		Amount:  1000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != "payment-1" {
		t.Errorf("expected ID payment-1, got %q", out.ID)
	}
	if out.Amount != 1000 {
		t.Errorf("expected amount 1000, got %d", out.Amount)
	}
	if out.Status != StatusPending {
		t.Errorf("expected status pending, got %q", out.Status)
	}
	if out.CreatedAt != 2000 {
		t.Errorf("expected created at 2000, got %d", out.CreatedAt)
	}
	if repo.saved == nil {
		t.Error("payment not saved")
	}
}

func TestCreatePayment_InvalidAmount(t *testing.T) {
	repo := &fakeRepo{}
	ids := &fixedIDGen{id: "payment-1"}
	clock := &fixedClock{now: 2000}
	uc := NewCreatePayment(repo, ids, clock, nil)

	_, err := uc.Execute(context.Background(), CreatePaymentInput{OrderID: "order-1", Amount: 0})
	if !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestCreatePayment_NegativeAmount(t *testing.T) {
	repo := &fakeRepo{}
	ids := &fixedIDGen{id: "payment-1"}
	clock := &fixedClock{now: 2000}
	uc := NewCreatePayment(repo, ids, clock, nil)

	_, err := uc.Execute(context.Background(), CreatePaymentInput{OrderID: "order-1", Amount: -100})
	if !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestGetPayment_Success(t *testing.T) {
	p := &Payment{ID: "payment-1", OrderID: "order-1", Amount: 1000, Status: StatusPending}
	repo := &fakeRepo{saved: p}
	uc := NewGetPayment(repo)

	out, err := uc.Execute(context.Background(), GetPaymentInput{ID: "payment-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != "payment-1" {
		t.Errorf("expected ID payment-1, got %q", out.ID)
	}
}

func TestGetPayment_NotFound(t *testing.T) {
	repo := &fakeRepo{}
	uc := NewGetPayment(repo)

	_, err := uc.Execute(context.Background(), GetPaymentInput{ID: "nonexistent"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetPaymentsByOrder(t *testing.T) {
	repo := &fakeRepo{}
	uc := NewGetPaymentsByOrder(repo)

	_, err := uc.Execute(context.Background(), GetPaymentsByOrderInput{OrderID: "order-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAllPayments(t *testing.T) {
	repo := &fakeRepo{}
	uc := NewGetAllPayments(repo)

	out, err := uc.Execute(context.Background(), GetAllPaymentsInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Payments == nil {
		t.Error("expected payments slice, got nil")
	}
}

func TestUpdatePayment_Success(t *testing.T) {
	repo := &fakeRepo{saved: &Payment{ID: "payment-1", OrderID: "order-1", Amount: 1000, Status: StatusPending}}
	uc := NewUpdatePayment(repo)

	status := StatusPaid
	out, err := uc.Execute(context.Background(), UpdatePaymentInput{
		ID:     "payment-1",
		Status: &status,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != StatusPaid {
		t.Errorf("expected status paid, got %q", out.Status)
	}
}

func TestUpdatePayment_InvalidStatus(t *testing.T) {
	repo := &fakeRepo{saved: &Payment{ID: "payment-1", OrderID: "order-1", Amount: 1000, Status: StatusPending}}
	uc := NewUpdatePayment(repo)

	status := Status("invalid")
	_, err := uc.Execute(context.Background(), UpdatePaymentInput{
		ID:     "payment-1",
		Status: &status,
	})
	if !errors.Is(err, ErrInvalidStatus) {
		t.Errorf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestUpdatePayment_NotFound(t *testing.T) {
	repo := &fakeRepo{}
	uc := NewUpdatePayment(repo)

	_, err := uc.Execute(context.Background(), UpdatePaymentInput{ID: "nonexistent"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestRefundPayment_Success(t *testing.T) {
	repo := &fakeRepo{saved: &Payment{ID: "payment-1", OrderID: "order-1", Amount: 1000, Status: StatusPaid}}
	uc := NewRefundPayment(repo, nil)

	out, err := uc.Execute(context.Background(), RefundPaymentInput{ID: "payment-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != StatusRefunded {
		t.Errorf("expected status refunded, got %q", out.Status)
	}
	if repo.saved == nil || repo.saved.Status != StatusRefunded {
		t.Error("payment not saved with refunded status")
	}
}

func TestRefundPayment_NotPaid(t *testing.T) {
	repo := &fakeRepo{saved: &Payment{ID: "payment-1", OrderID: "order-1", Amount: 1000, Status: StatusPending}}
	uc := NewRefundPayment(repo, nil)

	_, err := uc.Execute(context.Background(), RefundPaymentInput{ID: "payment-1"})
	if !errors.Is(err, ErrCannotRefund) {
		t.Errorf("expected ErrCannotRefund, got %v", err)
	}
}

func TestRefundPayment_Failed(t *testing.T) {
	repo := &fakeRepo{saved: &Payment{ID: "payment-1", OrderID: "order-1", Amount: 1000, Status: StatusFailed}}
	uc := NewRefundPayment(repo, nil)

	_, err := uc.Execute(context.Background(), RefundPaymentInput{ID: "payment-1"})
	if !errors.Is(err, ErrCannotRefund) {
		t.Errorf("expected ErrCannotRefund, got %v", err)
	}
}

func TestRefundPayment_AlreadyRefunded(t *testing.T) {
	repo := &fakeRepo{saved: &Payment{ID: "payment-1", OrderID: "order-1", Amount: 1000, Status: StatusRefunded}}
	uc := NewRefundPayment(repo, nil)

	_, err := uc.Execute(context.Background(), RefundPaymentInput{ID: "payment-1"})
	if !errors.Is(err, ErrCannotRefund) {
		t.Errorf("expected ErrCannotRefund, got %v", err)
	}
}

func TestRefundPayment_NotFound(t *testing.T) {
	repo := &fakeRepo{}
	uc := NewRefundPayment(repo, nil)

	_, err := uc.Execute(context.Background(), RefundPaymentInput{ID: "nonexistent"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}