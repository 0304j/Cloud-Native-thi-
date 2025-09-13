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

type CheckoutRequest struct {
	EventType    string       `json:"event_type"`
	OrderID      string       `json:"order_id"`
	UserID       string       `json:"user_id"`
	Items        []OrderItem  `json:"items"`
	TotalAmount  float64      `json:"total_amount"`
	Currency     string       `json:"currency"`
	Status       string       `json:"status"`
	Timestamp    string       `json:"timestamp"`
	
	// Restaurant delivery options
	OrderType    string       `json:"order_type"`              // "delivery" | "pickup"
	DeliveryInfo *DeliveryInfo `json:"delivery_info,omitempty"` // only for delivery orders
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
	
	// Validate order type
	if c.OrderType != "delivery" && c.OrderType != "pickup" {
		return fmt.Errorf("order_type must be 'delivery' or 'pickup'")
	}
	
	// Validate delivery info for delivery orders
	if c.OrderType == "delivery" {
		if c.DeliveryInfo == nil {
			return fmt.Errorf("delivery_info is required for delivery orders")
		}
		if err := c.DeliveryInfo.Validate(); err != nil {
			return fmt.Errorf("delivery_info validation failed: %w", err)
		}
	}
	
	return nil
}

func (d *DeliveryInfo) Validate() error {
	if d.CustomerName == "" {
		return fmt.Errorf("customer_name is required")
	}
	if d.CustomerPhone == "" {
		return fmt.Errorf("customer_phone is required")
	}
	if d.Street == "" {
		return fmt.Errorf("street is required")
	}
	if d.HouseNumber == "" {
		return fmt.Errorf("house_number is required")
	}
	if d.PostalCode == "" {
		return fmt.Errorf("postal_code is required")
	}
	if d.City == "" {
		return fmt.Errorf("city is required")
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
	order.OrderType = c.OrderType
	order.DeliveryInfo = c.DeliveryInfo
	return order
}
