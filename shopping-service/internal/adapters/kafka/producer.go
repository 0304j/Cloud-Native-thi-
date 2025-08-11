package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &KafkaProducer{writer: writer}
}

func (kp *KafkaProducer) SendMessage(value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msgBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = kp.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("product"),
		Value: msgBytes,
	})
	if err != nil {
		log.Println("Kafka send error:", err)
	}
	return err
}
