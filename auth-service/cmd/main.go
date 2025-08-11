package main

import (
	httpAdapter "auth-service/internal/adapters/http"
	mongoAdapter "auth-service/internal/adapters/mongo"
	"auth-service/internal/service"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	collection := os.Getenv("MONGO_COLLECTION")
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY"))

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbName)
	repo := mongoAdapter.NewUserRepository(db, collection)
	authService := service.NewAuthService(repo, jwtSecret, jwtExpiry)

	router := gin.Default()
	handler := httpAdapter.NewAuthHandler(authService)

	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	log.Println("Auth Service läuft auf Port 8081...")
	router.Run(":8081")
}
