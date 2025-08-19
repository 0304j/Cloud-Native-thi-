package ports

import (
	"context"
)

type OrderService interface {
	ProcessCheckout(ctx context.Context, checkoutData []byte) error
}
