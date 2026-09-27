package user

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, email, name string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	now := time.Now().Unix()
	u := &User{
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	return u, nil
}

func (s *Service) Get(ctx context.Context, id string) (*User, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetAll(ctx context.Context) ([]*User, error) {
	return s.repo.GetAll(ctx)
}

type UpdateInput struct {
	Email *string
	Name  *string
}

func (s *Service) Update(ctx context.Context, id string, in UpdateInput) (*User, error) {
	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Email != nil {
		if !isValidEmail(*in.Email) {
			return nil, ErrInvalidEmail
		}
		u.Email = *in.Email
	}
	if in.Name != nil {
		u.Name = *in.Name
	}
	u.UpdatedAt = time.Now().Unix()
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.Get(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}