package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

type Payment struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	UserID    uuid.UUID
	Provider  string
	Amount    float64
	Currency  string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPayment(orderID, userID uuid.UUID, provider string, amount float64, currency string) (*Payment, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	return &Payment{
		ID:        uuid.New(),
		OrderID:   orderID,
		UserID:    userID,
		Provider:  provider,
		Amount:    amount,
		Currency:  currency,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (p *Payment) MarkSuccess() error {
	if p.Status != StatusPending {
		return errors.New("only pending payments can be marked as success")
	}
	p.Status = StatusSuccess
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) MarkFailed() error {
	if p.Status != StatusPending {
		return errors.New("only pending payments can be marked as failed")
	}
	p.Status = StatusFailed
	p.UpdatedAt = time.Now()
	return nil
}
