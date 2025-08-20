package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"shopping-service/internal/adapters/http"
	"shopping-service/internal/adapters/http/middleware"
	"shopping-service/internal/adapters/kafka"
	mongoadapter "shopping-service/internal/adapters/mongo"
	"shopping-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 🔗 MongoDB Verbindung aus Environment Variable
	mongoURI := "mongodb://shopuser:shoppass@mongo:27017/shopping" // Default
	if envURI := os.Getenv("MONGO_URI"); envURI != "" {
		mongoURI = envURI
	}
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
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

	// CORS Middleware hinzufügen
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// CORS Preflight Handler für alle Routen
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

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
	
	// Health Check Endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "shopping-service"})
	})
	
	fmt.Println("🟢 Shopping Service läuft auf http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
