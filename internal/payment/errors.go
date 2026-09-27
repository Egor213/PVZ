package payment

import "errors"

var (
	ErrNotFound      = errors.New("payment not found")
	ErrInvalidAmount = errors.New("amount must be positive")
	ErrInvalidStatus = errors.New("invalid status")
	ErrCannotRefund  = errors.New("cannot refund payment")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}