package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"checkout-service/internal/domain/models"
	"checkout-service/internal/ports"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type PaymentEventConsumer struct {
	reader       *kafka.Reader
	orderService OrderDataProvider
}

type OrderDataProvider interface {
	GetOrderData(orderID string) *models.Order
}

func NewPaymentEventConsumer(orderService OrderDataProvider) *PaymentEventConsumer {
	return &PaymentEventConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   "payment-events",
			GroupID: "checkout-payment-group",
		}),
		orderService: orderService,
	}
}

func (c *PaymentEventConsumer) StartConsuming(ctx context.Context, eventPublisher ports.EventPublisher) error {
	log.Println("Starting to consume from 'payment-events' topic")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping payment event consumer")
			return c.reader.Close()
		default:
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Failed to read payment event: %v", err)
				continue
			}

			log.Printf("Received payment event: %s", string(message.Value))

			if err := c.processPaymentEvent(ctx, message.Value, eventPublisher); err != nil {
				log.Printf("Failed to process payment event: %v", err)
				continue
			}
		}
	}
}

func (c *PaymentEventConsumer) processPaymentEvent(ctx context.Context, eventData []byte, eventPublisher ports.EventPublisher) error {
	// Try to parse as PaymentSucceeded first
	var paymentSucceeded models.PaymentSucceededEvent
	if err := json.Unmarshal(eventData, &paymentSucceeded); err == nil && paymentSucceeded.OrderID != "" {
		return c.handlePaymentSucceeded(ctx, &paymentSucceeded, eventPublisher)
	}

	// Try to parse as PaymentFailed
	var paymentFailed models.PaymentFailedEvent
	if err := json.Unmarshal(eventData, &paymentFailed); err == nil && paymentFailed.OrderID != "" {
		return c.handlePaymentFailed(ctx, &paymentFailed, eventPublisher)
	}

	log.Printf("Unknown payment event format: %s", string(eventData))
	return nil
}

func (c *PaymentEventConsumer) handlePaymentSucceeded(ctx context.Context, event *models.PaymentSucceededEvent, eventPublisher ports.EventPublisher) error {
	log.Printf("Payment succeeded for order %s", event.OrderID)

	// Get original order data from store
	orderData := c.orderService.GetOrderData(event.OrderID)
	if orderData == nil {
		log.Printf("Warning: No order data found for order %s, using minimal data", event.OrderID)
		orderData = &models.Order{
			Items: []string{},
		}
	}

	var userID string
	if orderData != nil && orderData.UserID != uuid.Nil {
		userID = orderData.UserID.String()
	}

	// Use detailed items if available, otherwise fall back to product IDs
	var items []models.KitchenOrderItem
	if len(orderData.ItemDetails) > 0 {
		items = orderData.ItemDetails
	} else {
		// Fallback: convert product IDs to minimal KitchenOrderItems
		items = make([]models.KitchenOrderItem, len(orderData.Items))
		for i, productID := range orderData.Items {
			items[i] = models.KitchenOrderItem{
				ProductID:   productID,
				ProductName: "Unknown Product", // Kitchen will need to look this up
				Quantity:    1,                 // Default quantity
				Price:       0,                 // Unknown price
				PrepTime:    30,                // Default prep time
			}
		}
	}

	orderConfirmed := &models.OrderConfirmedEvent{
		EventType: "order_confirmed",
		OrderID:   event.OrderID,
		UserID:    userID,
		Items:     items,
		Timestamp: time.Now(),
	}

	if err := eventPublisher.PublishOrderConfirmed(ctx, orderConfirmed); err != nil {
		log.Printf("Failed to publish order confirmed event: %v", err)
		return err
	}

	log.Printf("Published order confirmed event for order %s with %d items", event.OrderID, len(items))
	return nil
}

func (c *PaymentEventConsumer) handlePaymentFailed(ctx context.Context, event *models.PaymentFailedEvent, eventPublisher ports.EventPublisher) error {
	log.Printf("Payment failed for order %s: %s", event.OrderID, event.Reason)

	orderCancelled := &models.OrderCancelledEvent{
		OrderID:   event.OrderID,
		UserID:    event.UserID,
		Reason:    "Payment failed: " + event.Reason,
		Timestamp: time.Now(),
	}

	if err := eventPublisher.PublishOrderCancelled(ctx, orderCancelled); err != nil {
		log.Printf("Failed to publish order cancelled event: %v", err)
		return err
	}

	log.Printf("Published order cancelled event for order %s", event.OrderID)
	return nil
}