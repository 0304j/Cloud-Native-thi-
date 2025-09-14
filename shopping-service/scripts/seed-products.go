package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"shopping-service/internal/adapters/mongo"
	"shopping-service/internal/domain"
	"shopping-service/internal/service"

	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("Seeding Analytica Restaurant products...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := "mongodb://shopuser:shoppass@mongo:27017/shopping"
	if envURI := os.Getenv("MONGO_URI"); envURI != "" {
		mongoURI = envURI
	}

	client, err := mongodriver.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	fmt.Println("Connected to MongoDB")
	db := client.Database("shopping")

	repo := mongo.NewMongoRepository(db)
	productService := service.NewProductService(repo)

	existingProducts, err := productService.GetAllProducts()
	if err != nil {
		log.Fatal("Failed to check existing products:", err)
	}

	if len(existingProducts) > 0 {
		fmt.Printf("Database already contains %d products. Skipping seed.\n", len(existingProducts))
		return
	}

	products := []domain.Product{
		{
			Name:        "Datenbrot",
			Description: "Hausgemachtes Brot mit Analytics-Butter und frischen Kräutern",
			Price:       6.50,
			UserID:      "seed-admin",
		},
		{
			Name:        "Algorithmus-Antipasti",
			Description: "Italienische Vorspeisen nach geheimem Algorithmus arrangiert",
			Price:       12.90,
			UserID:      "seed-admin",
		},
		{
			Name:        "Big Data Burger",
			Description: "Saftig gegrilltes Rindfleisch mit Analytics-Sauce und Hashbrown-Pommes",
			Price:       16.90,
			UserID:      "seed-admin",
		},
		{
			Name:        "Spaghetti Carbonara++",
			Description: "Klassische Carbonara mit optimiertem Algorithmus für perfekten Geschmack",
			Price:       14.50,
			UserID:      "seed-admin",
		},
		{
			Name:        "Cloud-Native Pizza",
			Description: "Pizza Margherita mit containerisierten Zutaten und microservice-Gewürzen",
			Price:       13.90,
			UserID:      "seed-admin",
		},
		{
			Name:        "Microservice Maultaschen",
			Description: "Handgemachte Maultaschen mit verteilter Füllung nach Microservice-Prinzip",
			Price:       11.90,
			UserID:      "seed-admin",
		},
		{
			Name:        "Container Curry",
			Description: "Exotisches Curry in perfekt orchestrierten Container-Portionen",
			Price:       15.50,
			UserID:      "seed-admin",
		},
		{
			Name:        "Kubernetes Käsespätzle",
			Description: "Schwäbische Spätzle mit automatisch skalierender Käsesauce",
			Price:       12.50,
			UserID:      "seed-admin",
		},
		{
			Name:        "Tiramisu 2.0",
			Description: "Traditionelles Tiramisu mit machine learning optimierter Mascarpone",
			Price:       7.50,
			UserID:      "seed-admin",
		},
		{
			Name:        "Docker Dampfnudeln",
			Description: "Fluffige Dampfnudeln aus dem Container mit süßer Vanilla-Sauce",
			Price:       8.90,
			UserID:      "seed-admin",
		},
	}

	successCount := 0
	for _, product := range products {
		err := productService.CreateProduct(&product)
		if err != nil {
			log.Printf("Failed to create product '%s': %v", product.Name, err)
		} else {
			fmt.Printf("Created: %s (EUR %.2f)\n", product.Name, product.Price)
			successCount++
		}
	}

	fmt.Printf("\nSeeding complete! Added %d/%d products to Analytica Restaurant menu.\n", successCount, len(products))
}
