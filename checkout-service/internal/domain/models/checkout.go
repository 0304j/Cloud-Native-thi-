package models

import (
	"fmt"

	"github.com/google/uuid"
)

type OrderItem struct {
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

type CheckoutRequest struct {
	EventType    string      `json:"event_type"`
	OrderID      string      `json:"order_id"`
	UserID       string      `json:"user_id"`
	Items        []OrderItem `json:"items"`
	ProductIDs   []string    `json:"product_ids"`
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

	// Convert OrderItems to simple product IDs for now
	productIDs := make([]string, len(c.Items))
	for i, item := range c.Items {
		productIDs[i] = item.ProductID
	}

	// Default payment provider - could be configurable
	paymentProvider := "stripe"

	return NewOrder(userID, productIDs, c.TotalAmount, paymentProvider)
}
