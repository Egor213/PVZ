package user

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type MemoryRepo struct {
	mu      sync.RWMutex
	items   map[string]*User
	byEmail map[string]*User
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		items:   make(map[string]*User),
		byEmail: make(map[string]*User),
	}
}

func (r *MemoryRepo) Save(ctx context.Context, u *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if existing, ok := r.byEmail[u.Email]; ok && existing.ID != u.ID {
		return ErrEmailExists
	}
	r.items[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}

func (r *MemoryRepo) Get(ctx context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

func (r *MemoryRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byEmail[email]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

func (r *MemoryRepo) GetAll(ctx context.Context) ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*User, 0, len(r.items))
	for _, u := range r.items {
		res = append(res, u)
	}
	return res, nil
}

func (r *MemoryRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.items[id]
	if !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	delete(r.byEmail, u.Email)
	return nil
}