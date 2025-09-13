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
	orderEventsWriter *kafka.Writer
	orderConfirmedWriter *kafka.Writer
	orderCancelledWriter *kafka.Writer
}

func NewEventPublisher() ports.EventPublisher {
	return &EventPublisher{
		orderEventsWriter: &kafka.Writer{
			Addr:     kafka.TCP("kafka:9092"),
			Topic:    "checkout-events",
			Balancer: &kafka.LeastBytes{},
		},
		orderConfirmedWriter: &kafka.Writer{
			Addr:     kafka.TCP("kafka:9092"),
			Topic:    "kitchen-events",
			Balancer: &kafka.LeastBytes{},
		},
		orderCancelledWriter: &kafka.Writer{
			Addr:     kafka.TCP("kafka:9092"),
			Topic:    "kitchen-events",
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

	if err := p.orderEventsWriter.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to publish OrderCreated event: %v", err)
		return err
	}

	log.Printf("Published OrderCreated event for order %s", event.OrderID)
	return nil
}

func (p *EventPublisher) PublishOrderConfirmed(ctx context.Context, event *models.OrderConfirmedEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.OrderID),
		Value: eventData,
	}

	if err := p.orderConfirmedWriter.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to publish OrderConfirmed event: %v", err)
		return err
	}

	log.Printf("Published OrderConfirmed event for order %s", event.OrderID)
	return nil
}

func (p *EventPublisher) PublishOrderCancelled(ctx context.Context, event *models.OrderCancelledEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.OrderID),
		Value: eventData,
	}

	if err := p.orderCancelledWriter.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to publish OrderCancelled event: %v", err)
		return err
	}

	log.Printf("Published OrderCancelled event for order %s", event.OrderID)
	return nil
}
