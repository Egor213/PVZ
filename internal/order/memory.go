package order

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type MemoryRepo struct {
	mu    sync.RWMutex
	items map[string]*Order
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{items: make(map[string]*Order)}
}

func (r *MemoryRepo) Save(ctx context.Context, o *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	r.items[o.ID] = o
	return nil
}

func (r *MemoryRepo) Get(ctx context.Context, id string) (*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return o, nil
}

func (r *MemoryRepo) GetAll(ctx context.Context) ([]*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*Order, 0, len(r.items))
	for _, o := range r.items {
		res = append(res, o)
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