package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"kitchen-service/internal/domain/models"
)

type KafkaPublisher struct {
	producer sarama.SyncProducer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync producer: %w", err)
	}

	return &KafkaPublisher{
		producer: producer,
	}, nil
}

func (kp *KafkaPublisher) PublishOrderStatusChanged(ctx context.Context, event *models.OrderStatusChangedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal order status event: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: "kitchen-events",
		Key:   sarama.StringEncoder(event.OrderID),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := kp.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to publish order status changed event: %w", err)
	}

	log.Printf("Order status event published to partition %d at offset %d", partition, offset)
	return nil
}

func (kp *KafkaPublisher) PublishKitchenNotification(ctx context.Context, event *models.KitchenNotificationEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal kitchen notification: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: "kitchen-events",
		Key:   sarama.StringEncoder(event.OrderID),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := kp.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to publish kitchen notification: %w", err)
	}

	log.Printf("Kitchen notification published to partition %d at offset %d", partition, offset)
	return nil
}

func (kp *KafkaPublisher) Close() error {
	return kp.producer.Close()
}