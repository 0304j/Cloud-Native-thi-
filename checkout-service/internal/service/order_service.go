package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"checkout-service/internal/domain/models"
	"checkout-service/internal/ports"
)

type OrderService struct {
	eventPublisher ports.EventPublisher
	orderStore     map[string]*models.Order // In-memory store for order data
	mutex          sync.RWMutex
}

func NewOrderService(publisher ports.EventPublisher) ports.OrderService {
	return &OrderService{
		eventPublisher: publisher,
		orderStore:     make(map[string]*models.Order),
		mutex:          sync.RWMutex{},
	}
}

func (s *OrderService) ProcessCheckout(ctx context.Context, checkoutData []byte) error {
	log.Printf("Processing checkout request: %s", string(checkoutData))

	var checkoutReq models.CheckoutRequest
	if err := json.Unmarshal(checkoutData, &checkoutReq); err != nil {
		log.Printf("Failed to parse checkout data: %v", err)
		return fmt.Errorf("invalid checkout data: %w", err)
	}

	if err := checkoutReq.Validate(); err != nil {
		log.Printf("Validation failed: %v", err)
		return err
	}

	order := checkoutReq.ToOrder()
	event := order.ToEvent()

	// Store order data for later use by payment callback
	s.mutex.Lock()
	s.orderStore[order.ID.String()] = order
	s.mutex.Unlock()

	if err := s.eventPublisher.PublishOrderCreated(ctx, event); err != nil {
		log.Printf("Failed to publish order created event: %v", err)
		return err
	}

	log.Printf("Successfully processed order %s with items: %v, total: %.2f",
		order.ID.String(), order.Items, order.TotalAmount)
	return nil
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
	log.Printf("Processing order %s with type: %s", order.ID.String(), order.OrderType)

	// Store order data
	s.mutex.Lock()
	s.orderStore[order.ID.String()] = order
	s.mutex.Unlock()

	// Create and publish event
	event := order.ToEvent()
	if err := s.eventPublisher.PublishOrderCreated(ctx, event); err != nil {
		log.Printf("Failed to publish order created event: %v", err)
		return err
	}

	log.Printf("Successfully processed order %s with type: %s, total: %.2f",
		order.ID.String(), order.OrderType, order.TotalAmount)
	return nil
}

func (s *OrderService) GetOrderData(orderID string) *models.Order {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.orderStore[orderID]
}
