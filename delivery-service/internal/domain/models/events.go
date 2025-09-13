package models

import "time"

// Incoming Events

type OrderCreatedEvent struct {
	OrderID         string             `json:"order_id"`
	UserID          string             `json:"user_id"`
	TotalAmount     float64            `json:"total_amount"`
	Currency        string             `json:"currency"`
	PaymentProvider string             `json:"payment_provider"`
	Items           []OrderItem        `json:"items"`
	CreatedAt       time.Time          `json:"created_at"`
	
	// Restaurant delivery info
	OrderType       string             `json:"order_type"`              // "delivery" | "pickup"
	DeliveryInfo    *DeliveryInfo      `json:"delivery_info,omitempty"` // only for delivery orders
}

type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	PrepTime    int     `json:"prep_time"`
}

type KitchenStatusChangedEvent struct {
	EventType     string      `json:"event_type"`
	OrderID       string      `json:"order_id"`
	Status        string      `json:"status"`        // "ready", "preparing", etc.
	EstimatedTime int         `json:"estimated_time,omitempty"`
	Timestamp     time.Time   `json:"timestamp"`
}

// Outgoing Events

type DeliveryAssignedEvent struct {
	EventType    string    `json:"event_type"`
	OrderID      string    `json:"order_id"`
	DriverID     string    `json:"driver_id"`
	CustomerInfo CustomerInfo `json:"customer_info"`
	Timestamp    time.Time `json:"timestamp"`
}

type DeliveryStatusEvent struct {
	EventType string         `json:"event_type"`
	OrderID   string         `json:"order_id"`
	Status    DeliveryStatus `json:"status"`
	Message   string         `json:"message"`
	Location  *GPSLocation   `json:"location,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

type PickupReadyEvent struct {
	EventType string    `json:"event_type"`
	OrderID   string    `json:"order_id"`
	Timestamp time.Time `json:"timestamp"`
}

// Supporting Types

type CustomerInfo struct {
	Name         string           `json:"name"`
	Phone        string           `json:"phone"`
	Address      DeliveryAddress  `json:"address"`
	Instructions string           `json:"instructions,omitempty"`
}

type DeliveryAddress struct {
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	PostalCode  string `json:"postal_code"`
	City        string `json:"city"`
	Floor       string `json:"floor,omitempty"`
}

type DeliveryStatus string

const (
	DeliveryStatusAssigned   DeliveryStatus = "assigned"
	DeliveryStatusPickedUp   DeliveryStatus = "picked_up"
	DeliveryStatusInTransit  DeliveryStatus = "in_transit"
	DeliveryStatusDelivered  DeliveryStatus = "delivered"
	DeliveryStatusCancelled  DeliveryStatus = "cancelled"
)

type GPSLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}