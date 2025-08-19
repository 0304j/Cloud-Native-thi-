package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"checkout-service/internal/domain/models"
	"checkout-service/internal/ports"
)

type OrderService struct {
	eventPublisher ports.EventPublisher
}

func NewOrderService(publisher ports.EventPublisher) ports.OrderService {
	return &OrderService{
		eventPublisher: publisher,
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

	if err := s.eventPublisher.PublishOrderCreated(ctx, event); err != nil {
		log.Printf("Failed to publish order created event: %v", err)
		return err
	}

	log.Printf("Successfully processed order %s with items: %v, total: %.2f",
		order.ID.String(), order.Items, order.TotalAmount)
	return nil
}
