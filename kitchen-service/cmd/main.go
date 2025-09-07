package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	httpHandler "kitchen-service/internal/adapters/http"
	"kitchen-service/internal/adapters/kafka"
	mongoRepo "kitchen-service/internal/adapters/mongo"
	"kitchen-service/internal/service"
)

func main() {
	log.Println("🍳 Starting Kitchen Service...")

	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	kafkaBrokers := []string{getEnv("KAFKA_BROKER", "localhost:9092")}
	port := getEnv("PORT", "8084")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(ctx)

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("✅ Connected to MongoDB")

	database := mongoClient.Database("kitchen_db")
	kitchenRepo := mongoRepo.NewKitchenRepository(database)

	publisher, err := kafka.NewKafkaPublisher(kafkaBrokers)
	if err != nil {
		log.Fatal("Failed to create Kafka publisher:", err)
	}
	defer publisher.Close()
	log.Println("✅ Kafka Publisher created")

	kitchenService := service.NewKitchenService(kitchenRepo, publisher)

	consumer, err := kafka.NewKafkaConsumer(
		kafkaBrokers,
		"kitchen-service-group",
		[]string{"kitchen-events"},
		kitchenService,
	)
	if err != nil {
		log.Fatal("Failed to create Kafka consumer:", err)
	}
	defer consumer.Close()
	log.Println("✅ Kafka Consumer created")

	go func() {
		log.Println("🎧 Starting Kafka consumer...")
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := kitchenService.ProcessOrderQueue(ctx); err != nil {
					log.Printf("Error processing order queue: %v", err)
				}
			}
		}
	}()

	router := mux.NewRouter()
	
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Kitchen Service is healthy! 🍳"))
	}).Methods("GET")

	kitchenHandler := httpHandler.NewKitchenHandler(kitchenService)
	kitchenHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("🚀 Kitchen Service running on port %s", port)
		log.Printf("📊 Dashboard available at: http://localhost:%s/kitchen/dashboard", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("🛑 Shutting down Kitchen Service...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("👋 Kitchen Service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}