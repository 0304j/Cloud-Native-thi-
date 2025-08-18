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
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id" binding:"required"`
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	Provider  string    `json:"provider" binding:"required,min=1,max=32"`
	Amount    float64   `json:"amount" binding:"required,gt=0"`
	Currency  string    `json:"currency" binding:"required,len=3"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
