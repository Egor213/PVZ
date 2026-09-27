package user

import (
	"context"
	"errors"
	"testing"
)

func TestService_Create(t *testing.T) {
	svc := NewService(NewMemoryRepo())
	u, err := svc.Create(context.Background(), "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", u.Email)
	}
	if u.Name != "Test User" {
		t.Errorf("expected name 'Test User', got %q", u.Name)
	}
	if u.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestService_Create_InvalidEmail(t *testing.T) {
	svc := NewService(NewMemoryRepo())
	_, err := svc.Create(context.Background(), "invalid-email", "Name")
	if !errors.Is(err, ErrInvalidEmail) {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestService_Create_DuplicateEmail(t *testing.T) {
	svc := NewService(NewMemoryRepo())
	svc.Create(context.Background(), "test@example.com", "User 1")
	_, err := svc.Create(context.Background(), "test@example.com", "User 2")
	if !errors.Is(err, ErrEmailExists) {
		t.Errorf("expected ErrEmailExists, got %v", err)
	}
}

func TestService_GetByEmail(t *testing.T) {
	svc := NewService(NewMemoryRepo())
	svc.Create(context.Background(), "test@example.com", "Test User")
	u, err := svc.GetByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", u.Email)
	}
}

func TestService_Update(t *testing.T) {
	svc := NewService(NewMemoryRepo())
	u, _ := svc.Create(context.Background(), "old@example.com", "Old Name")
	newName := "New Name"
	updated, err := svc.Update(context.Background(), u.ID, UpdateInput{Name: &newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "New Name" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}
}