package models

import (
	"time"
)

type DeliveryAggregate struct {
	OrderID          string           `json:"order_id"`
	
	// Event States
	OrderReceived    bool             `json:"order_received"`
	KitchenReady     bool             `json:"kitchen_ready"`
	DeliveryStarted  bool             `json:"delivery_started"`
	
	// Order Data
	UserID           string           `json:"user_id"`
	OrderType        string           `json:"order_type"`       // "delivery" | "pickup"
	TotalAmount      float64          `json:"total_amount"`
	Currency         string           `json:"currency"`
	
	// Delivery Data (only for delivery orders)
	DeliveryInfo     *DeliveryInfo    `json:"delivery_info,omitempty"`
	
	// Metadata
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	Status           AggregateStatus  `json:"status"`
	
	// Delivery Management
	DriverID         *string          `json:"driver_id,omitempty"`
}

type DeliveryInfo struct {
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	Street          string `json:"street"`
	HouseNumber     string `json:"house_number"`
	PostalCode      string `json:"postal_code"`
	City            string `json:"city"`
	Floor           string `json:"floor,omitempty"`
	Instructions    string `json:"instructions,omitempty"`
}

type AggregateStatus string

const (
	StatusWaitingForOrder   AggregateStatus = "waiting_for_order"
	StatusWaitingForKitchen AggregateStatus = "waiting_for_kitchen"
	StatusReadyToDeliver    AggregateStatus = "ready_to_deliver"
	StatusDelivering        AggregateStatus = "delivering"
	StatusCompleted         AggregateStatus = "completed"
	StatusPickupReady       AggregateStatus = "pickup_ready"   // for pickup orders
)

func (a *DeliveryAggregate) CanStartDelivery() bool {
	return a.OrderReceived && a.KitchenReady && !a.DeliveryStarted && a.OrderType == "delivery"
}

func (a *DeliveryAggregate) IsPickupReady() bool {
	return a.OrderReceived && a.KitchenReady && a.OrderType == "pickup"
}

func (a *DeliveryAggregate) UpdateStatus() {
	switch {
	case a.IsPickupReady():
		a.Status = StatusPickupReady
	case a.CanStartDelivery():
		a.Status = StatusReadyToDeliver
	case a.DeliveryStarted:
		a.Status = StatusDelivering
	case a.OrderReceived && !a.KitchenReady:
		a.Status = StatusWaitingForKitchen
	case !a.OrderReceived && a.KitchenReady:
		a.Status = StatusWaitingForOrder
	default:
		a.Status = StatusWaitingForOrder
	}
}