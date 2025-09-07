package kafka

import (
	"context"
	"log"

	"checkout-service/internal/ports"

	"github.com/segmentio/kafka-go"
)

type CheckoutConsumer struct {
	reader *kafka.Reader
}

func NewCheckoutConsumer() *CheckoutConsumer {
	return &CheckoutConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   "shopping-events",
			GroupID: "checkout-group",
		}),
	}
}

func (c *CheckoutConsumer) StartConsuming(ctx context.Context, orderService ports.OrderService) error {
	log.Println("Starting to consume from 'checkout' topic")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping consumer")
			return c.reader.Close()
		default:
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Failed to read message: %v", err)
				continue
			}

			if err := orderService.ProcessCheckout(ctx, message.Value); err != nil {
				log.Printf("Failed to process checkout: %v", err)
				continue
			}
		}
	}
}
