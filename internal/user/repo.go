package user

import "context"

type Repository interface {
	Save(ctx context.Context, u *User) error
	Get(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id string) error
}