package main

import (
	"log"
	"os"

	"payment-service/internal/adapters/http"
	"payment-service/internal/adapters/postgres"
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
	gormpostgres "gorm.io/driver/postgres" // Alias für gorm.io/driver/postgres
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	// DSN aus ENV (DATABASE_URL) oder Fallback
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Beispiel-Fallback; passe Werte an deine Umgebung an
		dsn = "host=postgres user=paymentuser password=paymentpass dbname=paymentdb port=5432 sslmode=disable TimeZone=UTC"
	}

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	repo := postgres.NewPaymentRepository(db)
	svc := service.NewService(repo)
	handler := &http.PaymentHandler{Service: svc}

	r.GET("/payments", handler.GetAllPayments)
	r.GET("/payments/:id", handler.GetPayment)
	r.POST("/payments", handler.CreatePayment)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
