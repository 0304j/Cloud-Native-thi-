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
	paymentConsumer := kafka.NewPaymentEventConsumer(orderService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start checkout consumer
	go func() {
		if err := consumer.StartConsuming(ctx, orderService); err != nil {
			log.Fatal("Checkout consumer failed:", err)
		}
	}()

	// Start payment event consumer
	go func() {
		if err := paymentConsumer.StartConsuming(ctx, publisher); err != nil {
			log.Fatal("Payment event consumer failed:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Checkout Service")
	cancel()
}
