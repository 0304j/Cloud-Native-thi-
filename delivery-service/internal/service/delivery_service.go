package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"delivery-service/internal/domain/models"
	"delivery-service/internal/ports"
)

type DeliveryService struct {
	repo      ports.AggregateRepository
	publisher ports.EventPublisher
}

func NewDeliveryService(repo ports.AggregateRepository, publisher ports.EventPublisher) ports.DeliveryService {
	return &DeliveryService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *DeliveryService) HandleOrderCreated(ctx context.Context, event *models.OrderCreatedEvent) error {
	log.Printf("Processing OrderCreated event for order %s, type: %s", event.OrderID, event.OrderType)

	// Get or create aggregate
	aggregate, err := s.repo.GetAggregate(ctx, event.OrderID)
	if err != nil {
		return fmt.Errorf("failed to get aggregate: %w", err)
	}

	if aggregate == nil {
		// Create new aggregate
		aggregate = &models.DeliveryAggregate{
			OrderID:   event.OrderID,
			Status:    models.StatusWaitingForKitchen,
			CreatedAt: time.Now(),
		}
	}

	// Apply Order Created Event
	aggregate.OrderReceived = true
	aggregate.UserID = event.UserID
	aggregate.OrderType = event.OrderType
	aggregate.TotalAmount = event.TotalAmount
	aggregate.Currency = event.Currency
	aggregate.DeliveryInfo = event.DeliveryInfo
	aggregate.UpdatedAt = time.Now()
	aggregate.UpdateStatus()

	// Save aggregate
	if err := s.repo.SaveAggregate(ctx, aggregate); err != nil {
		return fmt.Errorf("failed to save aggregate: %w", err)
	}

	log.Printf("Order %s marked as received, waiting for kitchen", event.OrderID)

	// Check if ready for next step
	return s.checkAndTriggerActions(ctx, aggregate)
}

func (s *DeliveryService) HandleKitchenStatusChanged(ctx context.Context, event *models.KitchenStatusChangedEvent) error {
	log.Printf("Processing KitchenStatusChanged event for order %s, status: %s", event.OrderID, event.Status)

	if event.Status != "ready" {
		log.Printf("Ignoring kitchen status %s for order %s", event.Status, event.OrderID)
		return nil
	}

	// Get aggregate
	aggregate, err := s.repo.GetAggregate(ctx, event.OrderID)
	if err != nil {
		return fmt.Errorf("failed to get aggregate: %w", err)
	}

	if aggregate == nil {
		// Create minimal aggregate - waiting for order event
		aggregate = &models.DeliveryAggregate{
			OrderID:   event.OrderID,
			Status:    models.StatusWaitingForOrder,
			CreatedAt: time.Now(),
		}
	}

	// Apply Kitchen Ready Event
	aggregate.KitchenReady = true
	aggregate.UpdatedAt = time.Now()
	aggregate.UpdateStatus()

	// Save aggregate
	if err := s.repo.SaveAggregate(ctx, aggregate); err != nil {
		return fmt.Errorf("failed to save aggregate: %w", err)
	}

	log.Printf("Order %s marked as kitchen ready", event.OrderID)

	// Check if ready for next step
	return s.checkAndTriggerActions(ctx, aggregate)
}

func (s *DeliveryService) checkAndTriggerActions(ctx context.Context, aggregate *models.DeliveryAggregate) error {
	switch {
	case aggregate.CanStartDelivery():
		return s.startDelivery(ctx, aggregate)
	case aggregate.IsPickupReady():
		return s.notifyPickupReady(ctx, aggregate)
	default:
		log.Printf("Order %s not ready yet - waiting for more events", aggregate.OrderID)
		return nil
	}
}

func (s *DeliveryService) startDelivery(ctx context.Context, aggregate *models.DeliveryAggregate) error {
	log.Printf("🚚 Starting delivery for order %s", aggregate.OrderID)

	// Mark as delivery started
	aggregate.DeliveryStarted = true
	aggregate.Status = models.StatusReadyToDeliver
	aggregate.UpdatedAt = time.Now()

	// TODO: Assign driver (simplified for now)
	driverID := "driver-001" // In real system: find available driver
	aggregate.DriverID = &driverID

	// Save updated aggregate
	if err := s.repo.SaveAggregate(ctx, aggregate); err != nil {
		return fmt.Errorf("failed to save aggregate after delivery start: %w", err)
	}

	// Create customer info from delivery info
	customerInfo := models.CustomerInfo{
		Name:  aggregate.DeliveryInfo.CustomerName,
		Phone: aggregate.DeliveryInfo.CustomerPhone,
		Address: models.DeliveryAddress{
			Street:      aggregate.DeliveryInfo.Street,
			HouseNumber: aggregate.DeliveryInfo.HouseNumber,
			PostalCode:  aggregate.DeliveryInfo.PostalCode,
			City:        aggregate.DeliveryInfo.City,
			Floor:       aggregate.DeliveryInfo.Floor,
		},
		Instructions: aggregate.DeliveryInfo.Instructions,
	}

	// Publish Delivery Assigned Event
	event := &models.DeliveryAssignedEvent{
		EventType:    "delivery_assigned",
		OrderID:      aggregate.OrderID,
		DriverID:     driverID,
		CustomerInfo: customerInfo,
		Timestamp:    time.Now(),
	}

	if err := s.publisher.PublishDeliveryAssigned(ctx, event); err != nil {
		log.Printf("Failed to publish delivery assigned event: %v", err)
		return err
	}

	log.Printf("✅ Delivery assigned for order %s to driver %s", aggregate.OrderID, driverID)
	
	// Start driver simulation for demo purposes
	go s.simulateDriverUpdates(context.Background(), aggregate.OrderID)
	
	return nil
}

func (s *DeliveryService) notifyPickupReady(ctx context.Context, aggregate *models.DeliveryAggregate) error {
	log.Printf("🏪 Order %s ready for pickup", aggregate.OrderID)

	// Update status
	aggregate.Status = models.StatusPickupReady
	aggregate.UpdatedAt = time.Now()

	// Save updated aggregate
	if err := s.repo.SaveAggregate(ctx, aggregate); err != nil {
		return fmt.Errorf("failed to save aggregate after pickup ready: %w", err)
	}

	// Publish Pickup Ready Event
	event := &models.PickupReadyEvent{
		EventType: "pickup_ready",
		OrderID:   aggregate.OrderID,
		Timestamp: time.Now(),
	}

	if err := s.publisher.PublishPickupReady(ctx, event); err != nil {
		log.Printf("Failed to publish pickup ready event: %v", err)
		return err
	}

	log.Printf("✅ Pickup ready notification sent for order %s", aggregate.OrderID)
	return nil
}

func (s *DeliveryService) GetOrderStatus(ctx context.Context, orderID string) (*models.DeliveryAggregate, error) {
	return s.repo.GetAggregate(ctx, orderID)
}

func (s *DeliveryService) UpdateDeliveryStatus(ctx context.Context, orderID string, status models.DeliveryStatus, message string) error {
	log.Printf("Updating delivery status for order %s to %s", orderID, status)

	// Get aggregate
	aggregate, err := s.repo.GetAggregate(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get aggregate: %w", err)
	}

	if aggregate == nil {
		return fmt.Errorf("order %s not found", orderID)
	}

	// Only update delivery orders
	if aggregate.OrderType != "delivery" {
		return fmt.Errorf("order %s is not a delivery order", orderID)
	}

	// Update aggregate with new status
	aggregate.UpdatedAt = time.Now()

	// Save updated aggregate
	if err := s.repo.SaveAggregate(ctx, aggregate); err != nil {
		return fmt.Errorf("failed to save aggregate: %w", err)
	}

	// Publish Delivery Status Event
	event := &models.DeliveryStatusEvent{
		EventType: "delivery_status_update",
		OrderID:   orderID,
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
	}

	if err := s.publisher.PublishDeliveryStatus(ctx, event); err != nil {
		log.Printf("Failed to publish delivery status event: %v", err)
		return err
	}

	log.Printf("✅ Delivery status updated for order %s: %s - %s", orderID, status, message)
	return nil
}

// simulateDriverUpdates simulates a driver updating delivery status for demo purposes
func (s *DeliveryService) simulateDriverUpdates(ctx context.Context, orderID string) {
	log.Printf("🤖 Starting driver simulation for order %s", orderID)
	
	// Step 1: Driver picks up the order (10-20 seconds delay)
	delay1 := time.Duration(10+rand.Intn(10)) * time.Second
	log.Printf("🤖 Driver will pick up order %s in %v", orderID, delay1)
	time.Sleep(delay1)
	
	err := s.UpdateDeliveryStatus(ctx, orderID, models.DeliveryStatusPickedUp, "Driver hat die Bestellung abgeholt")
	if err != nil {
		log.Printf("❌ Failed to update pickup status for order %s: %v", orderID, err)
		return
	}
	
	// Step 2: Driver starts transit (15-30 seconds delay)
	delay2 := time.Duration(15+rand.Intn(15)) * time.Second
	log.Printf("🤖 Driver will start transit for order %s in %v", orderID, delay2)
	time.Sleep(delay2)
	
	err = s.UpdateDeliveryStatus(ctx, orderID, models.DeliveryStatusInTransit, "Driver ist unterwegs zum Kunden")
	if err != nil {
		log.Printf("❌ Failed to update transit status for order %s: %v", orderID, err)
		return
	}
	
	// Step 3: Driver delivers the order (20-40 seconds delay)
	delay3 := time.Duration(20+rand.Intn(20)) * time.Second
	log.Printf("🤖 Driver will deliver order %s in %v", orderID, delay3)
	time.Sleep(delay3)
	
	err = s.UpdateDeliveryStatus(ctx, orderID, models.DeliveryStatusDelivered, "Bestellung wurde erfolgreich zugestellt")
	if err != nil {
		log.Printf("❌ Failed to update delivered status for order %s: %v", orderID, err)
		return
	}
	
	log.Printf("✅ Driver simulation completed for order %s", orderID)
}