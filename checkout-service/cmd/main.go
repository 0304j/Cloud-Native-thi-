package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "checkout",
		GroupID: "checkout-group",
	})

	fmt.Println("🟢 Checkout Service hört auf Kafka Topic: 'checkout'")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("❌ Kafka Fehler:", err)
		}
		fmt.Printf("📦 Empfangenes Produkt: %s\n", string(m.Value))
	}
}
