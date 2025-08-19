package models

import (
	"fmt"

	"github.com/google/uuid"
)

type CheckoutRequest struct {
	UserID          string   `json:"user_id,omitempty"`
	Items           []string `json:"items"`
	Total           float64  `json:"total"`
	Currency        string   `json:"currency,omitempty"`
	PaymentProvider string   `json:"payment_provider"`
}

func (c *CheckoutRequest) Validate() error {
	if len(c.Items) == 0 {
		return fmt.Errorf("items cannot be empty")
	}
	if c.Total <= 0 {
		return fmt.Errorf("total must be positive")
	}
	if c.PaymentProvider == "" {
		return fmt.Errorf("payment_provider is required")
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

	return NewOrder(userID, c.Items, c.Total, c.PaymentProvider)
}
