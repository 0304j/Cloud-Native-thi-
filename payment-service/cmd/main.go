package main

import (
	"context"
	"log"
	"os"

	"payment-service/internal/adapters/http"
	"payment-service/internal/adapters/postgres"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {

	r := gin.Default()

	conn, err := pgx.Connect(context.Background(), "postgres://root:rootpass@postgres:5432/payments")
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer conn.Close(context.Background())

	repo := &postgres.PaymentRepository{Conn: conn}
	handler := &http.PaymentHandler{Repo: repo}

	r.GET("/payments", handler.GetAllPayments)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	r.Run(":" + port)
}
