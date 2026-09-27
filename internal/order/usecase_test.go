package order

import (
	"context"
	"errors"
	"testing"
)

type fakeRepo struct {
	saved  *Order
	orders map[string]*Order
}

func (r *fakeRepo) Save(ctx context.Context, o *Order) error {
	r.saved = o
	if r.orders == nil {
		r.orders = make(map[string]*Order)
	}
	r.orders[o.ID] = o
	return nil
}

func (r *fakeRepo) Get(ctx context.Context, id string) (*Order, error) {
	if r.orders != nil {
		if o, ok := r.orders[id]; ok {
			return o, nil
		}
	}
	return nil, ErrNotFound
}

func (r *fakeRepo) GetAll(ctx context.Context) ([]*Order, error) {
	res := make([]*Order, 0, len(r.orders))
	for _, o := range r.orders {
		res = append(res, o)
	}
	return res, nil
}

func (r *fakeRepo) Delete(ctx context.Context, id string) error {
	return nil
}

type fixedIDGen struct{ id string }

func (g *fixedIDGen) New() string { return g.id }

type fixedClock struct{ now int64 }

func (c *fixedClock) Now() int64 { return c.now }

func TestCreateOrder_Success(t *testing.T) {
	repo := &fakeRepo{}
	ids := &fixedIDGen{id: "order-1"}
	clock := &fixedClock{now: 1000}
	uc := NewCreateOrder(repo, ids, clock, nil)

	out, err := uc.Execute(context.Background(), CreateOrderInput{
		Title:       "Test Order",
		Description: "Description",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != "order-1" {
		t.Errorf("expected ID order-1, got %q", out.ID)
	}
	if out.Title != "Test Order" {
		t.Errorf("expected title Test Order, got %q", out.Title)
	}
	if out.Status != StatusPending {
		t.Errorf("expected status pending, got %q", out.Status)
	}
	if out.CreatedAt != 1000 {
		t.Errorf("expected created at 1000, got %d", out.CreatedAt)
	}
	if repo.saved == nil {
		t.Error("order not saved")
	}
}

func TestCreateOrder_EmptyTitle(t *testing.T) {
	repo := &fakeRepo{}
	ids := &fixedIDGen{id: "order-1"}
	clock := &fixedClock{now: 1000}
	uc := NewCreateOrder(repo, ids, clock, nil)

	_, err := uc.Execute(context.Background(), CreateOrderInput{Title: ""})
	if !errors.Is(err, ErrEmptyTitle) {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestCreateOrder_TitleTooLong(t *testing.T) {
	repo := &fakeRepo{}
	ids := &fixedIDGen{id: "order-1"}
	clock := &fixedClock{now: 1000}
	uc := NewCreateOrder(repo, ids, clock, nil)

	longTitle := string(make([]byte, MaxTitleLength+1))
	_, err := uc.Execute(context.Background(), CreateOrderInput{Title: longTitle})
	if !errors.Is(err, ErrTitleTooLong) {
		t.Errorf("expected ErrTitleTooLong, got %v", err)
	}
}

func TestGetOrder_Success(t *testing.T) {
	o := &Order{ID: "order-1", Title: "Test", Status: StatusPending}
	repo := &fakeRepo{orders: map[string]*Order{"order-1": o}}
	uc := NewGetOrder(repo)

	out, err := uc.Execute(context.Background(), GetOrderInput{ID: "order-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != "order-1" {
		t.Errorf("expected ID order-1, got %q", out.ID)
	}
}

func TestGetOrder_NotFound(t *testing.T) {
	repo := &fakeRepo{orders: map[string]*Order{}}
	uc := NewGetOrder(repo)

	_, err := uc.Execute(context.Background(), GetOrderInput{ID: "nonexistent"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetAllOrders(t *testing.T) {
	o1 := &Order{ID: "order-1", Title: "Test 1", Status: StatusPending}
	o2 := &Order{ID: "order-2", Title: "Test 2", Status: StatusDone}
	repo := &fakeRepo{orders: map[string]*Order{"order-1": o1, "order-2": o2}}
	uc := NewGetAllOrders(repo)

	out, err := uc.Execute(context.Background(), GetAllOrdersInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(out.Orders))
	}
}

func TestUpdateOrder_Success(t *testing.T) {
	repo := &fakeRepo{saved: &Order{ID: "order-1", Title: "Original", Status: StatusPending}, orders: map[string]*Order{"order-1": {ID: "order-1", Title: "Original", Status: StatusPending}}}
	uc := NewUpdateOrder(repo)

	newTitle := "Updated"
	out, err := uc.Execute(context.Background(), UpdateOrderInput{
		ID:    "order-1",
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Title != "Updated" {
		t.Errorf("expected title Updated, got %q", out.Title)
	}
}

func TestUpdateOrder_Status(t *testing.T) {
	repo := &fakeRepo{saved: &Order{ID: "order-1", Title: "Test", Status: StatusPending}, orders: map[string]*Order{"order-1": {ID: "order-1", Title: "Test", Status: StatusPending}}}
	uc := NewUpdateOrder(repo)

	status := StatusInProgress
	out, err := uc.Execute(context.Background(), UpdateOrderInput{
		ID:     "order-1",
		Status: &status,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != StatusInProgress {
		t.Errorf("expected status in_progress, got %q", out.Status)
	}
}

func TestUpdateOrder_InvalidStatus(t *testing.T) {
	repo := &fakeRepo{saved: &Order{ID: "order-1", Title: "Test", Status: StatusPending}, orders: map[string]*Order{"order-1": {ID: "order-1", Title: "Test", Status: StatusPending}}}
	uc := NewUpdateOrder(repo)

	status := Status("invalid")
	_, err := uc.Execute(context.Background(), UpdateOrderInput{
		ID:     "order-1",
		Status: &status,
	})
	if !errors.Is(err, ErrInvalidStatus) {
		t.Errorf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestUpdateOrder_NotFound(t *testing.T) {
	repo := &fakeRepo{}
	uc := NewUpdateOrder(repo)

	_, err := uc.Execute(context.Background(), UpdateOrderInput{ID: "nonexistent"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCancelOrder_Success(t *testing.T) {
	repo := &fakeRepo{saved: &Order{ID: "order-1", Title: "Test", Status: StatusPending}, orders: map[string]*Order{"order-1": {ID: "order-1", Title: "Test", Status: StatusPending}}}
	uc := NewCancelOrder(repo, nil)

	out, err := uc.Execute(context.Background(), CancelOrderInput{ID: "order-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != StatusCancelled {
		t.Errorf("expected status cancelled, got %q", out.Status)
	}
	if repo.saved == nil || repo.saved.Status != StatusCancelled {
		t.Error("order not saved with cancelled status")
	}
}

func TestCancelOrder_AlreadyDone(t *testing.T) {
	repo := &fakeRepo{saved: &Order{ID: "order-1", Title: "Test", Status: StatusDone}, orders: map[string]*Order{"order-1": {ID: "order-1", Title: "Test", Status: StatusDone}}}
	uc := NewCancelOrder(repo, nil)

	_, err := uc.Execute(context.Background(), CancelOrderInput{ID: "order-1"})
	if !errors.Is(err, ErrCannotCancel) {
		t.Errorf("expected ErrCannotCancel, got %v", err)
	}
}

func TestCancelOrder_AlreadyCancelled(t *testing.T) {
	repo := &fakeRepo{saved: &Order{ID: "order-1", Title: "Test", Status: StatusCancelled}, orders: map[string]*Order{"order-1": {ID: "order-1", Title: "Test", Status: StatusCancelled}}}
	uc := NewCancelOrder(repo, nil)

	_, err := uc.Execute(context.Background(), CancelOrderInput{ID: "order-1"})
	if !errors.Is(err, ErrCannotCancel) {
		t.Errorf("expected ErrCannotCancel, got %v", err)
	}
}

func TestCancelOrder_NotFound(t *testing.T) {
	repo := &fakeRepo{}
	uc := NewCancelOrder(repo, nil)

	_, err := uc.Execute(context.Background(), CancelOrderInput{ID: "nonexistent"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}