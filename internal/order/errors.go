package order

import "errors"

var (
	ErrNotFound      = errors.New("order not found")
	ErrEmptyTitle    = errors.New("title cannot be empty")
	ErrTitleTooLong  = errors.New("title exceeds maximum length")
	ErrInvalidStatus = errors.New("invalid status")
	ErrCannotCancel  = errors.New("cannot cancel order")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}