package ports

import (
	"context"
	"checkout-service/internal/domain/models"
)

type OrderService interface {
	ProcessCheckout(ctx context.Context, checkoutData []byte) error
	GetOrderData(orderID string) *models.Order
}
