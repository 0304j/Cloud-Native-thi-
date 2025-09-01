package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	StatusReceived         OrderStatus = "received"
	StatusPreparing        OrderStatus = "preparing"
	StatusReady            OrderStatus = "ready"
	StatusPickedUpByDriver OrderStatus = "picked_up_by_driver"
	StatusCancelled        OrderStatus = "cancelled"
)

type KitchenOrder struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderID        string             `json:"order_id" bson:"order_id"`
	CustomerID     string             `json:"customer_id" bson:"customer_id"`
	Items          []OrderItem        `json:"items" bson:"items"`
	Status         OrderStatus        `json:"status" bson:"status"`
	EstimatedTime  int                `json:"estimated_time" bson:"estimated_time"` // in seconds
	StartedAt      *time.Time         `json:"started_at,omitempty" bson:"started_at,omitempty"`
	CompletedAt    *time.Time         `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type OrderItem struct {
	ProductID   string  `json:"product_id" bson:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Price       float64 `json:"price" bson:"price"`
	PrepTime    int     `json:"prep_time" bson:"prep_time"` // preparation time in seconds
}

type KitchenStats struct {
	TotalOrders            int `json:"total_orders"`
	OrdersReceived         int `json:"orders_received"`
	OrdersPreparing        int `json:"orders_preparing"`
	OrdersReady            int `json:"orders_ready"`
	OrdersPickedUpByDriver int `json:"orders_picked_up_by_driver"`
	AverageWaitTime        int `json:"average_wait_time"` // in seconds
}