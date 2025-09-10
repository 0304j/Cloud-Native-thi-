package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	httpAdapter "checkout-service/internal/adapters/http"
	"checkout-service/internal/adapters/kafka"
	"checkout-service/internal/service"
)

func main() {
	log.Println("Starting Checkout Service with hexagonal architecture")

	publisher := kafka.NewEventPublisher()
	consumer := kafka.NewCheckoutConsumer()
	
	orderService := service.NewOrderService(publisher)
	paymentConsumer := kafka.NewPaymentEventConsumer(orderService)
	
	// HTTP Server
	httpServer := httpAdapter.NewServer(orderService)
	
	// Get port from environment or default to 8082
	port := 8082
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil {
			port = p
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start HTTP server
	go func() {
		log.Printf("Starting HTTP server on port %d", port)
		if err := httpServer.Start(port); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server failed:", err)
		}
	}()

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
	
	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	
	cancel()
}
