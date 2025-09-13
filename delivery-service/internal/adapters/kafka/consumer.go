package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"delivery-service/internal/domain/models"
	"delivery-service/internal/ports"

	"github.com/segmentio/kafka-go"
)

type EventConsumer struct {
	checkoutReader *kafka.Reader
	kitchenReader  *kafka.Reader
	deliveryService ports.DeliveryService
}

func NewEventConsumer(brokers []string, deliveryService ports.DeliveryService) *EventConsumer {
	return &EventConsumer{
		checkoutReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       "checkout-events",
			GroupID:     "delivery-service-checkout",
			StartOffset: kafka.LastOffset,
		}),
		kitchenReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       "kitchen-events", 
			GroupID:     "delivery-service-kitchen",
			StartOffset: kafka.LastOffset,
		}),
		deliveryService: deliveryService,
	}
}

func (c *EventConsumer) StartConsuming(ctx context.Context) error {
	// Start checkout events consumer
	go func() {
		if err := c.consumeCheckoutEvents(ctx); err != nil {
			log.Printf("Checkout events consumer error: %v", err)
		}
	}()

	// Start kitchen events consumer
	go func() {
		if err := c.consumeKitchenEvents(ctx); err != nil {
			log.Printf("Kitchen events consumer error: %v", err)
		}
	}()

	return nil
}

func (c *EventConsumer) consumeCheckoutEvents(ctx context.Context) error {
	log.Println("Starting to consume checkout events...")
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.checkoutReader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Failed to fetch checkout message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if err := c.handleCheckoutMessage(ctx, msg); err != nil {
				log.Printf("Failed to handle checkout message: %v", err)
			}

			if err := c.checkoutReader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Failed to commit checkout message: %v", err)
			}
		}
	}
}

func (c *EventConsumer) consumeKitchenEvents(ctx context.Context) error {
	log.Println("Starting to consume kitchen events...")
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.kitchenReader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Failed to fetch kitchen message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if err := c.handleKitchenMessage(ctx, msg); err != nil {
				log.Printf("Failed to handle kitchen message: %v", err)
			}

			if err := c.kitchenReader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Failed to commit kitchen message: %v", err)
			}
		}
	}
}

func (c *EventConsumer) handleCheckoutMessage(ctx context.Context, msg kafka.Message) error {
	log.Printf("Received checkout event: %s", string(msg.Value))

	var event models.OrderCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	// Only process delivery and pickup orders
	if event.OrderType != "delivery" && event.OrderType != "pickup" {
		log.Printf("Ignoring order %s with type %s", event.OrderID, event.OrderType)
		return nil
	}

	return c.deliveryService.HandleOrderCreated(ctx, &event)
}

func (c *EventConsumer) handleKitchenMessage(ctx context.Context, msg kafka.Message) error {
	log.Printf("Received kitchen event: %s", string(msg.Value))

	// Try to parse as KitchenStatusChangedEvent
	var statusEvent models.KitchenStatusChangedEvent
	if err := json.Unmarshal(msg.Value, &statusEvent); err != nil {
		// If it doesn't match, try other event types
		log.Printf("Failed to parse as status event, trying other formats: %v", err)
		return nil
	}

	return c.deliveryService.HandleKitchenStatusChanged(ctx, &statusEvent)
}

func (c *EventConsumer) Close() error {
	if err := c.checkoutReader.Close(); err != nil {
		log.Printf("Error closing checkout reader: %v", err)
	}
	if err := c.kitchenReader.Close(); err != nil {
		log.Printf("Error closing kitchen reader: %v", err)
	}
	return nil
}