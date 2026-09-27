package payment

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type MemoryRepo struct {
	mu    sync.RWMutex
	items map[string]*Payment
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{items: make(map[string]*Payment)}
}

func (r *MemoryRepo) Save(ctx context.Context, p *Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	r.items[p.ID] = p
	return nil
}

func (r *MemoryRepo) Get(ctx context.Context, id string) (*Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (r *MemoryRepo) GetByOrderID(ctx context.Context, orderID string) ([]*Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*Payment, 0)
	for _, p := range r.items {
		if p.OrderID == orderID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (r *MemoryRepo) GetAll(ctx context.Context) ([]*Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*Payment, 0, len(r.items))
	for _, p := range r.items {
		res = append(res, p)
	}
	return res, nil
}

func (r *MemoryRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}