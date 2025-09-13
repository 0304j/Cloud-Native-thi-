package ports

import (
	"context"
	"delivery-service/internal/domain/models"
)

type AggregateRepository interface {
	SaveAggregate(ctx context.Context, aggregate *models.DeliveryAggregate) error
	GetAggregate(ctx context.Context, orderID string) (*models.DeliveryAggregate, error)
	DeleteAggregate(ctx context.Context, orderID string) error
	GetReadyToDeliverOrders(ctx context.Context) ([]*models.DeliveryAggregate, error)
}

type EventPublisher interface {
	PublishDeliveryAssigned(ctx context.Context, event *models.DeliveryAssignedEvent) error
	PublishDeliveryStatus(ctx context.Context, event *models.DeliveryStatusEvent) error
	PublishPickupReady(ctx context.Context, event *models.PickupReadyEvent) error
	Close() error
}

type DeliveryService interface {
	HandleOrderCreated(ctx context.Context, event *models.OrderCreatedEvent) error
	HandleKitchenStatusChanged(ctx context.Context, event *models.KitchenStatusChangedEvent) error
	GetOrderStatus(ctx context.Context, orderID string) (*models.DeliveryAggregate, error)
	UpdateDeliveryStatus(ctx context.Context, orderID string, status models.DeliveryStatus, message string) error
}