package mongo

import (
	"context"
	"shopping-service/internal/domain"
	"shopping-service/internal/ports"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) ports.ProductRepository {
	return &MongoRepository{
		collection: db.Collection("products"),
	}
}

func (m *MongoRepository) Create(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := m.collection.InsertOne(ctx, product)
	if err != nil {
		log.Printf("InsertOne Error: %v", err)
		return err
	}
	log.Printf("Inserted document ID: %v", result.InsertedID)
	return nil
}

func (m *MongoRepository) FindAll() ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("❌ Fehler bei Find():", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	for cursor.Next(ctx) {
		var p domain.Product
		if err := cursor.Decode(&p); err != nil {
			log.Println("❌ Fehler beim Dekodieren:", err)
			return nil, err
		}
		products = append(products, p)
	}

	if err := cursor.Err(); err != nil {
		log.Println("❌ Cursor-Fehler:", err)
		return nil, err
	}

	log.Printf("✅ %d Produkte gefunden", len(products))
	return products, nil
}

func (m *MongoRepository) FindByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert string ID to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("❌ Invalid ObjectID format: %s", id)
		return nil, err
	}

	var product domain.Product
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("❌ Product with ID %s not found", id)
			return nil, err
		}
		log.Printf("❌ Error finding product by ID %s: %v", id, err)
		return nil, err
	}

	log.Printf("✅ Product found: %s - %.2f", product.Name, product.Price)
	return &product, nil
}
