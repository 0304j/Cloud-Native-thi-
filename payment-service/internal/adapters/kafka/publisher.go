package kafka

import (
	"context"
	"encoding/json"
	"log"

	"payment-service/internal/domain/models"
	"payment-service/internal/ports"

	"github.com/segmentio/kafka-go"
)

type EventPublisher struct {
	writer *kafka.Writer
}

func NewEventPublisher() ports.EventPublisher {
	return &EventPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP("kafka:9092"),
			Topic:    "payment-events",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *EventPublisher) PublishPaymentSucceeded(ctx context.Context, event *models.PaymentSucceededEvent) error {
	return p.publishEvent(ctx, "PaymentSucceeded", event.OrderID, event)
}

func (p *EventPublisher) PublishPaymentFailed(ctx context.Context, event *models.PaymentFailedEvent) error {
	return p.publishEvent(ctx, "PaymentFailed", event.OrderID, event)
}

func (p *EventPublisher) publishEvent(ctx context.Context, eventType, orderID string, event interface{}) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(orderID),
		Value: eventData,
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to publish %s event: %v", eventType, err)
		return err
	}

	log.Printf("Published %s event for order %s", eventType, orderID)
	return nil
}
