package user

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt int64
	UpdatedAt int64
}

func NewUser(email, name string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	now := time.Now().Unix()
	return &User{
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) UpdateEmail(email string) error {
	if !isValidEmail(email) {
		return ErrInvalidEmail
	}
	u.Email = email
	u.UpdatedAt = time.Now().Unix()
	return nil
}

func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now().Unix()
}