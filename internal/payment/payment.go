package payment

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusPaid     Status = "paid"
	StatusFailed   Status = "failed"
	StatusRefunded Status = "refunded"
)

func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusPaid, StatusFailed, StatusRefunded:
		return true
	}
	return false
}

type Payment struct {
	ID        string
	OrderID   string
	Amount    int64
	Status    Status
	CreatedAt int64
	UpdatedAt int64
}

func NewPayment(orderID string, amount int64) (*Payment, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	now := time.Now().Unix()
	return &Payment{
		OrderID:   orderID,
		Amount:    amount,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (p *Payment) MarkPaid() error {
	if p.Status != StatusPending {
		return ErrInvalidStatus
	}
	p.Status = StatusPaid
	p.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *Payment) MarkFailed() error {
	if p.Status != StatusPending {
		return ErrInvalidStatus
	}
	p.Status = StatusFailed
	p.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *Payment) Refund() error {
	if p.Status != StatusPaid {
		return ErrCannotRefund
	}
	p.Status = StatusRefunded
	p.UpdatedAt = time.Now().Unix()
	return nil
}