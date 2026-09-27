package order

import "time"

const MaxTitleLength = 200

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusDone, StatusCancelled:
		return true
	}
	return false
}

type Order struct {
	ID          string
	Title       string
	Description string
	Status      Status
	CreatedAt   int64
	UpdatedAt   int64
}

func NewOrder(title, description string) (*Order, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if len(title) > MaxTitleLength {
		return nil, ErrTitleTooLong
	}
	now := time.Now().Unix()
	return &Order{
		Title:       title,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (o *Order) UpdateTitle(title string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if len(title) > MaxTitleLength {
		return ErrTitleTooLong
	}
	o.Title = title
	o.UpdatedAt = time.Now().Unix()
	return nil
}

func (o *Order) UpdateDescription(description string) {
	o.Description = description
	o.UpdatedAt = time.Now().Unix()
}

func (o *Order) UpdateStatus(status Status) error {
	if !status.Valid() {
		return ErrInvalidStatus
	}
	o.Status = status
	o.UpdatedAt = time.Now().Unix()
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == StatusDone {
		return ErrCannotCancel
	}
	if o.Status == StatusCancelled {
		return ErrCannotCancel
	}
	o.Status = StatusCancelled
	o.UpdatedAt = time.Now().Unix()
	return nil
}

func (o *Order) IsDone() bool {
	return o.Status == StatusDone
}

func (o *Order) IsCancelled() bool {
	return o.Status == StatusCancelled
}