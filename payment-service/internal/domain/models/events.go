package models

import "time"

type OrderCreatedEvent struct {
	OrderID         string    `json:"order_id"`
	UserID          string    `json:"user_id"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	PaymentProvider string    `json:"payment_provider"`
	Items           []string  `json:"items"`
	CreatedAt       time.Time `json:"created_at"`
}

type PaymentSucceededEvent struct {
	OrderID   string    `json:"order_id"`
	PaymentID string    `json:"payment_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

type PaymentFailedEvent struct {
	OrderID   string    `json:"order_id"`
	Reason    string    `json:"reason"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}
