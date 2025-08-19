package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	PaymentProvider string    `json:"payment_provider"`
	Items           []string  `json:"items"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type OrderCreatedEvent struct {
	OrderID         string    `json:"order_id"`
	UserID          string    `json:"user_id"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	PaymentProvider string    `json:"payment_provider"`
	Items           []string  `json:"items"`
	CreatedAt       time.Time `json:"created_at"`
}

func NewOrder(userID uuid.UUID, items []string, totalAmount float64, paymentProvider string) *Order {
	return &Order{
		ID:              uuid.New(),
		UserID:          userID,
		TotalAmount:     totalAmount,
		Currency:        "EUR",
		PaymentProvider: paymentProvider,
		Items:           items,
		Status:          "created",
		CreatedAt:       time.Now(),
	}
}

func (o *Order) ToEvent() *OrderCreatedEvent {
	return &OrderCreatedEvent{
		OrderID:         o.ID.String(),
		UserID:          o.UserID.String(),
		TotalAmount:     o.TotalAmount,
		Currency:        o.Currency,
		PaymentProvider: o.PaymentProvider,
		Items:           o.Items,
		CreatedAt:       o.CreatedAt,
	}
}
