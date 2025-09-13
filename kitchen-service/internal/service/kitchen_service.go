package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"kitchen-service/internal/domain/models"
	"kitchen-service/internal/ports"
)

type KitchenService struct {
	repo      ports.KitchenRepository
	publisher ports.EventPublisher
}

func NewKitchenService(repo ports.KitchenRepository, publisher ports.EventPublisher) *KitchenService {
	return &KitchenService{
		repo:      repo,
		publisher: publisher,
	}
}

func (ks *KitchenService) ReceiveOrder(ctx context.Context, order *models.KitchenOrder) error {
	order.Status = models.StatusReceived
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	if err := ks.repo.SaveOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to save received order: %w", err)
	}

	event := &models.OrderStatusChangedEvent{
		EventType:     "order_received_in_kitchen",
		OrderID:       order.OrderID,
		Status:        models.StatusReceived,
		EstimatedTime: order.EstimatedTime,
		Timestamp:     time.Now(),
	}

	if err := ks.publisher.PublishOrderStatusChanged(ctx, event); err != nil {
		log.Printf("Failed to publish order received event: %v", err)
	}

	notification := &models.KitchenNotificationEvent{
		EventType: "kitchen_notification",
		OrderID:   order.OrderID,
		Status:    models.StatusReceived,
		Message:   fmt.Sprintf("🍽️ Neue Bestellung erhalten! Geschätzte Zubereitungszeit: %d Sekunden", order.EstimatedTime),
		Timestamp: time.Now(),
	}

	if err := ks.publisher.PublishKitchenNotification(ctx, notification); err != nil {
		log.Printf("Failed to publish kitchen notification: %v", err)
	}

	log.Printf("🍳 Kitchen received order %s with %d items", order.OrderID, len(order.Items))
	return nil
}

func (ks *KitchenService) StartPreparation(ctx context.Context, orderID string) error {
	order, err := ks.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status != models.StatusReceived {
		return fmt.Errorf("order %s is not in received status, current: %s", orderID, order.Status)
	}

	now := time.Now()
	order.Status = models.StatusPreparing
	order.StartedAt = &now
	order.UpdatedAt = now

	if err := ks.repo.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	event := &models.OrderStatusChangedEvent{
		EventType:     "order_preparation_started",
		OrderID:       order.OrderID,
		Status:        models.StatusPreparing,
		EstimatedTime: order.EstimatedTime,
		Timestamp:     time.Now(),
	}

	if err := ks.publisher.PublishOrderStatusChanged(ctx, event); err != nil {
		log.Printf("Failed to publish preparation started event: %v", err)
	}

	notification := &models.KitchenNotificationEvent{
		EventType: "kitchen_notification",
		OrderID:   order.OrderID,
		Status:    models.StatusPreparing,
		Message:   "👨‍🍳 Die Zubereitung hat begonnen! Der Koch ist am Werk...",
		Timestamp: time.Now(),
	}

	if err := ks.publisher.PublishKitchenNotification(ctx, notification); err != nil {
		log.Printf("Failed to publish kitchen notification: %v", err)
	}

	go ks.simulatePreparation(ctx, order)

	log.Printf("🔥 Started preparing order %s", orderID)
	return nil
}

func (ks *KitchenService) simulatePreparation(ctx context.Context, order *models.KitchenOrder) {
	estimatedSeconds := order.EstimatedTime
	actualSeconds := ks.calculateActualPrepTime(estimatedSeconds)
	
	log.Printf("🕐 Simulating preparation for order %s: estimated %d sec, actual %d sec", 
		order.OrderID, estimatedSeconds, actualSeconds)

	preparationSteps := []string{
		"🥘 Zutaten werden vorbereitet...",
		"🔥 Kochen hat begonnen...",
		"👨‍🍳 Chef fügt Gewürze hinzu...",
		"⏰ Fast fertig...",
	}

	stepDuration := time.Duration(actualSeconds) * time.Second / time.Duration(len(preparationSteps))

	for i, step := range preparationSteps {
		select {
		case <-ctx.Done():
			return
		case <-time.After(stepDuration):
			notification := &models.KitchenNotificationEvent{
				EventType: "kitchen_notification",
				OrderID:   order.OrderID,
				Status:    models.StatusPreparing,
				Message:   fmt.Sprintf("%s (Schritt %d/%d)", step, i+1, len(preparationSteps)),
				Timestamp: time.Now(),
			}

			if err := ks.publisher.PublishKitchenNotification(ctx, notification); err != nil {
				log.Printf("Failed to publish preparation step notification: %v", err)
			}
		}
	}

	if err := ks.CompleteOrder(ctx, order.OrderID); err != nil {
		log.Printf("Failed to complete order automatically: %v", err)
	}
}

func (ks *KitchenService) calculateActualPrepTime(prepTimeSeconds int) int {
	// Use prep time directly as seconds
	variance := 0.2
	minTime := float64(prepTimeSeconds) * (1.0 - variance)
	maxTime := float64(prepTimeSeconds) * (1.0 + variance)
	
	actualTime := minTime + rand.Float64()*(maxTime-minTime)
	
	if actualTime < 5 {
		actualTime = 5 // minimum 5 seconds
	}
	
	return int(actualTime)
}

func (ks *KitchenService) CompleteOrder(ctx context.Context, orderID string) error {
	order, err := ks.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status != models.StatusPreparing {
		return fmt.Errorf("order %s is not being prepared, current: %s", orderID, order.Status)
	}

	now := time.Now()
	order.Status = models.StatusReady
	order.CompletedAt = &now
	order.UpdatedAt = now

	if err := ks.repo.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	event := &models.OrderStatusChangedEvent{
		EventType: "order_ready",
		OrderID:   order.OrderID,
		Status:    models.StatusReady,
		Timestamp: time.Now(),
	}

	if err := ks.publisher.PublishOrderStatusChanged(ctx, event); err != nil {
		log.Printf("Failed to publish order ready event: %v", err)
	}

	notification := &models.KitchenNotificationEvent{
		EventType: "kitchen_notification",
		OrderID:   order.OrderID,
		Status:    models.StatusReady,
		Message:   "✅ Bestellung ist fertig! Bereit zur Abholung 🍽️",
		Timestamp: time.Now(),
	}

	if err := ks.publisher.PublishKitchenNotification(ctx, notification); err != nil {
		log.Printf("Failed to publish kitchen notification: %v", err)
	}

	log.Printf("✅ Order %s is ready!", orderID)
	return nil
}

func (ks *KitchenService) MarkPickedUpByDriver(ctx context.Context, orderID string) error {
	order, err := ks.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status != models.StatusReady {
		return fmt.Errorf("order %s is not ready for pickup, current: %s", orderID, order.Status)
	}

	order.Status = models.StatusPickedUpByDriver
	order.UpdatedAt = time.Now()

	if err := ks.repo.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	event := &models.OrderStatusChangedEvent{
		EventType: "order_picked_up_by_driver",
		OrderID:   order.OrderID,
		Status:    models.StatusPickedUpByDriver,
		Timestamp: time.Now(),
	}

	if err := ks.publisher.PublishOrderStatusChanged(ctx, event); err != nil {
		log.Printf("Failed to publish order pickup event: %v", err)
	}

	notification := &models.KitchenNotificationEvent{
		EventType: "kitchen_notification",
		OrderID:   order.OrderID,
		Status:    models.StatusPickedUpByDriver,
		Message:   "🚗 Bestellung wurde vom Fahrer abgeholt! Auf dem Weg zum Kunden! 📦",
		Timestamp: time.Now(),
	}

	if err := ks.publisher.PublishKitchenNotification(ctx, notification); err != nil {
		log.Printf("Failed to publish kitchen notification: %v", err)
	}

	log.Printf("🚗 Order %s picked up by driver!", orderID)
	return nil
}

func (ks *KitchenService) CancelOrder(ctx context.Context, orderID string) error {
	order, err := ks.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status == models.StatusPickedUpByDriver {
		return fmt.Errorf("cannot cancel order %s - already picked up by driver", orderID)
	}

	order.Status = models.StatusCancelled
	order.UpdatedAt = time.Now()

	if err := ks.repo.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	event := &models.OrderStatusChangedEvent{
		EventType: "order_cancelled_in_kitchen",
		OrderID:   order.OrderID,
		Status:    models.StatusCancelled,
		Timestamp: time.Now(),
	}

	if err := ks.publisher.PublishOrderStatusChanged(ctx, event); err != nil {
		log.Printf("Failed to publish order cancelled event: %v", err)
	}

	notification := &models.KitchenNotificationEvent{
		EventType: "kitchen_notification",
		OrderID:   order.OrderID,
		Status:    models.StatusCancelled,
		Message:   "❌ Bestellung wurde storniert",
		Timestamp: time.Now(),
	}

	if err := ks.publisher.PublishKitchenNotification(ctx, notification); err != nil {
		log.Printf("Failed to publish kitchen notification: %v", err)
	}

	log.Printf("❌ Order %s cancelled", orderID)
	return nil
}

func (ks *KitchenService) GetOrderStatus(ctx context.Context, orderID string) (*models.KitchenOrder, error) {
	return ks.repo.GetOrderByID(ctx, orderID)
}

func (ks *KitchenService) GetAllOrders(ctx context.Context) ([]*models.KitchenOrder, error) {
	return ks.repo.GetAllOrders(ctx)
}

func (ks *KitchenService) GetKitchenStats(ctx context.Context) (*models.KitchenStats, error) {
	return ks.repo.GetStats(ctx)
}

func (ks *KitchenService) ProcessOrderQueue(ctx context.Context) error {
	receivedOrders, err := ks.repo.GetOrdersByStatus(ctx, models.StatusReceived)
	if err != nil {
		return fmt.Errorf("failed to get received orders: %w", err)
	}

	for _, order := range receivedOrders {
		if err := ks.StartPreparation(ctx, order.OrderID); err != nil {
			log.Printf("Failed to start preparation for order %s: %v", order.OrderID, err)
		}
		
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * 2):
		}
	}

	return nil
}