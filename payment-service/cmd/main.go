package main

import (
	"context"
	"log"
	"os"

	"payment-service/internal/adapters/http"
	"payment-service/internal/adapters/postgres"
	"payment-service/internal/service"

	"os/signal"
	"payment-service/internal/adapters/kafka"
	"syscall"

	"github.com/gin-gonic/gin"
	gormpostgres "gorm.io/driver/postgres" // Alias für gorm.io/driver/postgres
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=postgres user=paymentuser password=paymentpass dbname=paymentdb port=5432 sslmode=disable TimeZone=UTC"
	}

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	if err := db.AutoMigrate(&postgres.PaymentEntity{}); err != nil {
		log.Fatal("DB migration failed:", err)
	}

	repo := postgres.NewPaymentRepository(db)
	eventPublisher := kafka.NewEventPublisher()
	svc := service.NewService(repo, eventPublisher)
	handler := &http.PaymentHandler{Service: svc}

	consumer := kafka.NewOrderEventConsumer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.StartConsuming(ctx, svc); err != nil {
			log.Fatal("Order event consumer failed:", err)
		}
	}()

	r.GET("/payments", handler.GetAllPayments)
	r.GET("/payments/:id", handler.GetPayment)
	r.POST("/payments", handler.CreatePayment)
	r.PUT("/payments/:id/status", handler.UpdatePaymentStatus)
	r.DELETE("/payments/:id", handler.DeletePayment)
	r.GET("/payments/search", handler.GetPaymentsByStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(":" + port); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("Shutting down Payment Service")
	cancel()
}
