package service

import (
	"context"
	"log"

	"checkout-service/internal/domain/models"
	"checkout-service/internal/ports"

	"github.com/google/uuid"
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

	//TODO Parse checkoutData
	order := models.NewOrder(
		uuid.New(),                           // userID
		[]string{"Pizza Margherita", "Cola"}, // items
		29.99,                                // totalAmount
		"stripe",                             // paymentProvider
	)

	event := order.ToEvent()

	if err := s.eventPublisher.PublishOrderCreated(ctx, event); err != nil {
		log.Printf("Failed to publish order created event: %v", err)
		return err
	}

	log.Printf("Successfully processed order %s", order.ID.String())
	return nil
}
