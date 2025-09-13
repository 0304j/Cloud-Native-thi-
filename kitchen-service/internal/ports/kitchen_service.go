package ports

import (
	"context"
	"kitchen-service/internal/domain/models"
)

type KitchenService interface {
	ReceiveOrder(ctx context.Context, order *models.KitchenOrder) error
	StartPreparation(ctx context.Context, orderID string) error
	CompleteOrder(ctx context.Context, orderID string) error
	MarkPickedUpByDriver(ctx context.Context, orderID string) error
	CancelOrder(ctx context.Context, orderID string) error
	GetOrderStatus(ctx context.Context, orderID string) (*models.KitchenOrder, error)
	GetAllOrders(ctx context.Context) ([]*models.KitchenOrder, error)
	GetKitchenStats(ctx context.Context) (*models.KitchenStats, error)
	ProcessOrderQueue(ctx context.Context) error
}

type KitchenRepository interface {
	SaveOrder(ctx context.Context, order *models.KitchenOrder) error
	UpdateOrder(ctx context.Context, order *models.KitchenOrder) error
	GetOrderByID(ctx context.Context, orderID string) (*models.KitchenOrder, error)
	GetOrdersByStatus(ctx context.Context, status models.OrderStatus) ([]*models.KitchenOrder, error)
	GetAllOrders(ctx context.Context) ([]*models.KitchenOrder, error)
	DeleteOrder(ctx context.Context, orderID string) error
	GetStats(ctx context.Context) (*models.KitchenStats, error)
}

type EventPublisher interface {
	PublishOrderStatusChanged(ctx context.Context, event *models.OrderStatusChangedEvent) error
	PublishKitchenNotification(ctx context.Context, event *models.KitchenNotificationEvent) error
}