package models

import (
	"fmt"

	"github.com/google/uuid"
)

type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

type CheckoutRequest struct {
	EventType    string      `json:"event_type"`
	OrderID      string      `json:"order_id"`
	UserID       string      `json:"user_id"`
	Items        []OrderItem `json:"items"`
	TotalAmount  float64     `json:"total_amount"`
	Currency     string      `json:"currency"`
	Status       string      `json:"status"`
	Timestamp    string      `json:"timestamp"`
}

func (c *CheckoutRequest) Validate() error {
	if len(c.Items) == 0 {
		return fmt.Errorf("items cannot be empty")
	}
	if c.TotalAmount <= 0 {
		return fmt.Errorf("total amount must be positive")
	}
	if c.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	return nil
}

func (c *CheckoutRequest) ToOrder() *Order {
	userID := uuid.New()
	if c.UserID != "" {
		if parsed, err := uuid.Parse(c.UserID); err == nil {
			userID = parsed
		}
	}

	currency := c.Currency
	if currency == "" {
		currency = "EUR"
	}

	// Convert OrderItems to simple product IDs for backward compatibility
	productIDs := make([]string, len(c.Items))
	for i, item := range c.Items {
		productIDs[i] = item.ProductID
	}

	// Convert CheckoutRequest.OrderItems to KitchenOrderItems
	orderItems := make([]KitchenOrderItem, len(c.Items))
	for i, item := range c.Items {
		// Use product name from the order item if available
		productName := item.ProductName
		if productName == "" {
			productName = "Unknown Product"
		}
		
		// Debug log to see what we're getting
		fmt.Printf("Processing item %d: ProductID=%s, ProductName='%s', Quantity=%d\n", 
			i, item.ProductID, item.ProductName, item.Quantity)

		orderItems[i] = KitchenOrderItem{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Price:       item.UnitPrice,
			PrepTime:    30, // Default prep time - could be configurable
		}
	}

	// Default payment provider - could be configurable
	paymentProvider := "stripe"

	order := NewOrder(userID, productIDs, c.TotalAmount, paymentProvider)
	order.ItemDetails = orderItems
	return order
}
