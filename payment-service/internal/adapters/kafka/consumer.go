package kafka

import (
	"context"
	"encoding/json"
	"log"

	"payment-service/internal/domain/models"
	"payment-service/internal/ports"

	"github.com/segmentio/kafka-go"
)

type OrderEventConsumer struct {
	reader *kafka.Reader
}

func NewOrderEventConsumer() *OrderEventConsumer {
	return &OrderEventConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   "order-events",
			GroupID: "payment-service-group",
		}),
	}
}

func (c *OrderEventConsumer) StartConsuming(ctx context.Context, paymentService ports.PaymentService) error {
	log.Println("Starting to consume from 'order-events' topic")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping order event consumer")
			return c.reader.Close()
		default:
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Failed to read order event: %v", err)
				continue
			}

			var orderEvent models.OrderCreatedEvent
			if err := json.Unmarshal(message.Value, &orderEvent); err != nil {
				log.Printf("Failed to unmarshal order event: %v", err)
				continue
			}

			payment, err := paymentService.ProcessOrderPayment(ctx, orderEvent)
			if err != nil {
				log.Printf("Failed to process payment for order %s: %v", orderEvent.OrderID, err)
				continue
			}

			log.Printf("Payment %s processed for order %s", payment.ID.String(), orderEvent.OrderID)
		}
	}
}
