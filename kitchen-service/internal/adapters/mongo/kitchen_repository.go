package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"kitchen-service/internal/domain/models"
)

type KitchenRepository struct {
	collection *mongo.Collection
}

func NewKitchenRepository(database *mongo.Database) *KitchenRepository {
	return &KitchenRepository{
		collection: database.Collection("kitchen_orders"),
	}
}

func (kr *KitchenRepository) SaveOrder(ctx context.Context, order *models.KitchenOrder) error {
	if order.ID.IsZero() {
		order.ID = primitive.NewObjectID()
	}
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	_, err := kr.collection.InsertOne(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to save kitchen order: %w", err)
	}

	return nil
}

func (kr *KitchenRepository) UpdateOrder(ctx context.Context, order *models.KitchenOrder) error {
	order.UpdatedAt = time.Now()

	filter := bson.M{"order_id": order.OrderID}
	update := bson.M{"$set": order}

	result, err := kr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update kitchen order: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("kitchen order with ID %s not found", order.OrderID)
	}

	return nil
}

func (kr *KitchenRepository) GetOrderByID(ctx context.Context, orderID string) (*models.KitchenOrder, error) {
	var order models.KitchenOrder
	filter := bson.M{"order_id": orderID}

	err := kr.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("kitchen order with ID %s not found", orderID)
		}
		return nil, fmt.Errorf("failed to get kitchen order: %w", err)
	}

	return &order, nil
}

func (kr *KitchenRepository) GetOrdersByStatus(ctx context.Context, status models.OrderStatus) ([]*models.KitchenOrder, error) {
	filter := bson.M{"status": status}
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := kr.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find orders by status: %w", err)
	}
	defer cursor.Close(ctx)

	var orders []*models.KitchenOrder
	for cursor.Next(ctx) {
		var order models.KitchenOrder
		if err := cursor.Decode(&order); err != nil {
			return nil, fmt.Errorf("failed to decode kitchen order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return orders, nil
}

func (kr *KitchenRepository) GetAllOrders(ctx context.Context) ([]*models.KitchenOrder, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	cursor, err := kr.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find all orders: %w", err)
	}
	defer cursor.Close(ctx)

	var orders []*models.KitchenOrder
	for cursor.Next(ctx) {
		var order models.KitchenOrder
		if err := cursor.Decode(&order); err != nil {
			return nil, fmt.Errorf("failed to decode kitchen order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return orders, nil
}

func (kr *KitchenRepository) DeleteOrder(ctx context.Context, orderID string) error {
	filter := bson.M{"order_id": orderID}

	result, err := kr.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete kitchen order: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("kitchen order with ID %s not found", orderID)
	}

	return nil
}

func (kr *KitchenRepository) GetStats(ctx context.Context) (*models.KitchenStats, error) {
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": "$status",
				"count": bson.M{"$sum": 1},
				"avgTime": bson.M{"$avg": "$estimated_time"},
			},
		},
	}

	cursor, err := kr.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate stats: %w", err)
	}
	defer cursor.Close(ctx)

	stats := &models.KitchenStats{}
	statusCounts := make(map[models.OrderStatus]int)
	totalTime := 0.0
	totalOrders := 0

	for cursor.Next(ctx) {
		var result struct {
			ID      models.OrderStatus `bson:"_id"`
			Count   int                `bson:"count"`
			AvgTime float64            `bson:"avgTime"`
		}

		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode stats result: %w", err)
		}

		statusCounts[result.ID] = result.Count
		totalOrders += result.Count
		totalTime += result.AvgTime * float64(result.Count)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	stats.TotalOrders = totalOrders
	stats.OrdersReceived = statusCounts[models.StatusReceived]
	stats.OrdersPreparing = statusCounts[models.StatusPreparing]
	stats.OrdersReady = statusCounts[models.StatusReady]
	stats.OrdersPickedUpByDriver = statusCounts[models.StatusPickedUpByDriver]

	if totalOrders > 0 {
		stats.AverageWaitTime = int(totalTime / float64(totalOrders))
	}

	return stats, nil
}