package ports

import (
	"context"
	"payment-service/internal/domain/models"
)

type EventPublisher interface {
	PublishPaymentSucceeded(ctx context.Context, event *models.PaymentSucceededEvent) error
	PublishPaymentFailed(ctx context.Context, event *models.PaymentFailedEvent) error
}
