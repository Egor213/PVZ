package user

import (
	"errors"
	"strings"
)

var (
	ErrNotFound     = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
	ErrInvalidEmail = errors.New("invalid email")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsEmailExists(err error) bool {
	return errors.Is(err, ErrEmailExists)
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && len(email) > 3
}