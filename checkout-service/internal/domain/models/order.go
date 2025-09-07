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
	Items           []string  `json:"items"`          // Keep for backward compatibility
	ItemDetails     []KitchenOrderItem `json:"item_details"` // Store detailed item information
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type OrderCreatedEvent struct {
	OrderID         string             `json:"order_id"`
	UserID          string             `json:"user_id"`
	TotalAmount     float64            `json:"total_amount"`
	Currency        string             `json:"currency"`
	PaymentProvider string             `json:"payment_provider"`
	Items           []KitchenOrderItem `json:"items"`
	CreatedAt       time.Time          `json:"created_at"`
}

type OrderConfirmedEvent struct {
	EventType string              `json:"event_type"`
	OrderID   string              `json:"order_id"`
	UserID    string              `json:"user_id"`
	Items     []KitchenOrderItem  `json:"items"`
	Timestamp time.Time           `json:"timestamp"`
}

// OrderItem for kitchen events - more detailed than checkout OrderItem
type KitchenOrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	PrepTime    int     `json:"prep_time"`
}

type OrderCancelledEvent struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

type PaymentSucceededEvent struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	PaymentID string  `json:"payment_id"`
}

type PaymentFailedEvent struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Reason  string `json:"reason"`
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
		Items:           o.ItemDetails,
		CreatedAt:       o.CreatedAt,
	}
}
