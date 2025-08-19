package ports

import (
	"checkout-service/internal/domain/models"
	"context"
)

type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, event *models.OrderCreatedEvent) error
}
