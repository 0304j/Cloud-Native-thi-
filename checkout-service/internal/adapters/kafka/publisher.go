package kafka

import (
	"context"
	"encoding/json"
	"log"

	"checkout-service/internal/domain/models"
	"checkout-service/internal/ports"

	"github.com/segmentio/kafka-go"
)

type EventPublisher struct {
	writer *kafka.Writer
}

func NewEventPublisher() ports.EventPublisher {
	return &EventPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP("kafka:9092"),
			Topic:    "order-events",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *EventPublisher) PublishOrderCreated(ctx context.Context, event *models.OrderCreatedEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.OrderID),
		Value: eventData,
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to publish OrderCreated event: %v", err)
		return err
	}

	log.Printf("Published OrderCreated event for order %s", event.OrderID)
	return nil
}
