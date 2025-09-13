package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"delivery-service/internal/adapters/kafka"
	"delivery-service/internal/adapters/redis"
	"delivery-service/internal/service"

	redisClient "github.com/redis/go-redis/v9"
)

func main() {
	log.Println("🚚 Starting Delivery Service with Event Aggregation")

	// Redis setup
	redisAddr := getEnv("REDIS_ADDR", "redis:6379")
	rdb := redisClient.NewClient(&redisClient.Options{
		Addr: redisAddr,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")

	// Dependencies
	aggregateRepo := redis.NewAggregateRepository(rdb)
	eventPublisher := kafka.NewEventPublisher([]string{"kafka:9092"})
	deliveryService := service.NewDeliveryService(aggregateRepo, eventPublisher)

	// Kafka consumer
	eventConsumer := kafka.NewEventConsumer([]string{"kafka:9092"}, deliveryService)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming events
	if err := eventConsumer.StartConsuming(ctx); err != nil {
		log.Fatalf("Failed to start consuming events: %v", err)
	}
	log.Println("✅ Started consuming events from checkout-events and kitchen-events")

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Delivery Service...")

	// Graceful shutdown
	cancel()
	if err := eventConsumer.Close(); err != nil {
		log.Printf("Error closing event consumer: %v", err)
	}
	if err := eventPublisher.Close(); err != nil {
		log.Printf("Error closing event publisher: %v", err)
	}
	if err := rdb.Close(); err != nil {
		log.Printf("Error closing Redis client: %v", err)
	}

	log.Println("✅ Delivery Service stopped gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}