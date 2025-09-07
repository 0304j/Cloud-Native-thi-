package models

import "time"

type OrderReceivedEvent struct {
	EventType  string       `json:"event_type"`
	OrderID    string       `json:"order_id"`
	CustomerID string       `json:"customer_id"`
	UserID     string       `json:"user_id"`     // Support user_id from checkout service
	Items      []OrderItem  `json:"items"`
	Timestamp  time.Time    `json:"timestamp"`
}

type OrderStatusChangedEvent struct {
	EventType     string      `json:"event_type"`
	OrderID       string      `json:"order_id"`
	Status        OrderStatus `json:"status"`
	EstimatedTime int         `json:"estimated_time,omitempty"`
	Timestamp     time.Time   `json:"timestamp"`
}

type KitchenNotificationEvent struct {
	EventType string      `json:"event_type"`
	OrderID   string      `json:"order_id"`
	Status    OrderStatus `json:"status"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
}