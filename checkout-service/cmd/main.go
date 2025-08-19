package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"checkout-service/internal/adapters/kafka"
	"checkout-service/internal/service"
)

func main() {
	log.Println("Starting Checkout Service with hexagonal architecture")

	publisher := kafka.NewEventPublisher()
	consumer := kafka.NewCheckoutConsumer()

	orderService := service.NewOrderService(publisher)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.StartConsuming(ctx, orderService); err != nil {
			log.Fatal("Consumer failed:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Checkout Service")
	cancel()
}
