package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shopping-service/internal/adapters/http"
	"shopping-service/internal/adapters/http/middleware"
	"shopping-service/internal/adapters/kafka"
	mongoadapter "shopping-service/internal/adapters/mongo"
	"shopping-service/internal/service"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 🔗 MongoDB Verbindung
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://shopuser:shoppass@mongo:27017/shopping"))
	if err != nil {
		log.Fatal("MongoDB Verbindung fehlgeschlagen:", err)
	}
	// Ping MongoDB um Verbindung zu prüfen
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB Ping fehlgeschlagen:", err)
	}

	fmt.Println("✅ MongoDB Verbindung erfolgreich!")
	db := client.Database("shopping")

	// 💡 Layer zusammensetzen
	repo := mongoadapter.NewMongoRepository(db)
	productService := service.NewProductService(repo)

	// 🌐 HTTP starten
	r := gin.Default()

	// Kafka Producer initialisieren
	kafkaProducer := kafka.NewKafkaProducer("kafka:9092", "checkout")

	http.NewProductHandler(r, productService, kafkaProducer)
	userGroup := r.Group("/")
	userGroup.Use(middleware.JWTMiddleware())
	{
		cartRepo := mongoadapter.NewCartRepo(db)
		cartSvc := service.NewCartService(cartRepo)
		http.NewCartHandler(userGroup, cartSvc, kafkaProducer)
	}
	fmt.Println("🟢 Shopping Service läuft auf http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
