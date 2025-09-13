package kafka

import (
	"context"
	"encoding/json"
	"log"

	"delivery-service/internal/domain/models"
	"delivery-service/internal/ports"

	"github.com/segmentio/kafka-go"
)

type EventPublisher struct {
	writer *kafka.Writer
}

func NewEventPublisher(brokers []string) ports.EventPublisher {
	return &EventPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "delivery-events",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *EventPublisher) PublishDeliveryAssigned(ctx context.Context, event *models.DeliveryAssignedEvent) error {
	return p.publishEvent(ctx, event, event.OrderID)
}

func (p *EventPublisher) PublishDeliveryStatus(ctx context.Context, event *models.DeliveryStatusEvent) error {
	return p.publishEvent(ctx, event, event.OrderID)
}

func (p *EventPublisher) PublishPickupReady(ctx context.Context, event *models.PickupReadyEvent) error {
	return p.publishEvent(ctx, event, event.OrderID)
}

func (p *EventPublisher) publishEvent(ctx context.Context, event interface{}, orderID string) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(orderID),
		Value: eventData,
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to publish event: %v", err)
		return err
	}

	log.Printf("Published event to delivery-events: %s", string(eventData))
	return nil
}

func (p *EventPublisher) Close() error {
	return p.writer.Close()
}