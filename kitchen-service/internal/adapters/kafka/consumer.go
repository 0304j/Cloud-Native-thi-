package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"kitchen-service/internal/domain/models"
	"kitchen-service/internal/ports"
)

type KafkaConsumer struct {
	consumer       sarama.ConsumerGroup
	kitchenService ports.KitchenService
	topics         []string
}

func NewKafkaConsumer(brokers []string, groupID string, topics []string, kitchenService ports.KitchenService) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &KafkaConsumer{
		consumer:       consumer,
		kitchenService: kitchenService,
		topics:         topics,
	}, nil
}

func (kc *KafkaConsumer) Start(ctx context.Context) error {
	handler := &ConsumerGroupHandler{
		kitchenService: kc.kitchenService,
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Kitchen consumer context cancelled")
			return ctx.Err()
		default:
			if err := kc.consumer.Consume(ctx, kc.topics, handler); err != nil {
				log.Printf("Error consuming messages: %v", err)
				time.Sleep(time.Second * 5)
			}
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}

type ConsumerGroupHandler struct {
	kitchenService ports.KitchenService
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Kitchen consumer group setup")
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Kitchen consumer group cleanup")
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			log.Printf("Received message from topic %s: %s", message.Topic, string(message.Value))

			if err := h.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *ConsumerGroupHandler) processMessage(message *sarama.ConsumerMessage) error {
	ctx := context.Background()

	switch message.Topic {
	case "kitchen-events":
		return h.handleKitchenEvent(ctx, message.Value)
	default:
		log.Printf("Unknown topic: %s", message.Topic)
		return nil
	}
}

func (h *ConsumerGroupHandler) handleKitchenEvent(ctx context.Context, data []byte) error {
	// Parse event to determine type
	var baseEvent struct {
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(data, &baseEvent); err != nil {
		log.Printf("Failed to parse event type: %v", err)
		return err
	}

	// Only handle order confirmations - ignore our own status events!
	switch baseEvent.EventType {
	case "order_confirmed":
		return h.handleOrderConfirmed(ctx, data)
	case "order_received_in_kitchen", "order_preparation_started", "order_ready", 
		 "order_picked_up_by_driver", "order_cancelled_in_kitchen", "kitchen_notification":
		// Ignore our own status events to prevent infinite loop
		log.Printf("Ignoring own event type: %s", baseEvent.EventType)
		return nil
	default:
		log.Printf("Unknown kitchen event type: %s", baseEvent.EventType)
		return nil
	}
}

func (h *ConsumerGroupHandler) handleOrderConfirmed(ctx context.Context, data []byte) error {
	var orderEvent models.OrderReceivedEvent
	if err := json.Unmarshal(data, &orderEvent); err != nil {
		return fmt.Errorf("failed to unmarshal order event: %w", err)
	}

	// Use user_id as customer_id if customer_id is not set
	customerID := orderEvent.CustomerID
	if customerID == "" {
		customerID = orderEvent.UserID
	}

	kitchenOrder := &models.KitchenOrder{
		OrderID:       orderEvent.OrderID,
		CustomerID:    customerID,
		Items:         orderEvent.Items,
		Status:        models.StatusReceived,
		EstimatedTime: calculateEstimatedTime(orderEvent.Items),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.kitchenService.ReceiveOrder(ctx, kitchenOrder); err != nil {
		return fmt.Errorf("failed to receive order in kitchen: %w", err)
	}

	log.Printf("Order %s received in kitchen", orderEvent.OrderID)
	return nil
}

func (h *ConsumerGroupHandler) handleOrderCancelled(ctx context.Context, data []byte) error {
	var cancelEvent struct {
		OrderID string `json:"order_id"`
	}
	
	if err := json.Unmarshal(data, &cancelEvent); err != nil {
		return fmt.Errorf("failed to unmarshal cancel event: %w", err)
	}

	if err := h.kitchenService.CancelOrder(ctx, cancelEvent.OrderID); err != nil {
		return fmt.Errorf("failed to cancel order in kitchen: %w", err)
	}

	log.Printf("Order %s cancelled in kitchen", cancelEvent.OrderID)
	return nil
}

func calculateEstimatedTime(items []models.OrderItem) int {
	totalTime := 0
	for _, item := range items {
		if item.PrepTime > 0 {
			totalTime += item.PrepTime * item.Quantity
		} else {
			totalTime += 30 * item.Quantity // default 30 seconds per item
		}
	}
	
	if totalTime < 10 {
		totalTime = 10 // minimum 10 seconds
	}
	if totalTime > 300 {
		totalTime = 300 // maximum 5 minutes
	}
	
	return totalTime
}