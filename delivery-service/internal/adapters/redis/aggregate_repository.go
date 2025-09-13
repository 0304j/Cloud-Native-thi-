package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"delivery-service/internal/domain/models"
	"delivery-service/internal/ports"

	"github.com/redis/go-redis/v9"
)

type AggregateRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewAggregateRepository(redisClient *redis.Client) ports.AggregateRepository {
	return &AggregateRepository{
		client: redisClient,
		ttl:    24 * time.Hour, // 24h TTL for aggregates
	}
}

func (r *AggregateRepository) SaveAggregate(ctx context.Context, aggregate *models.DeliveryAggregate) error {
	key := r.getKey(aggregate.OrderID)
	
	data, err := json.Marshal(aggregate)
	if err != nil {
		return fmt.Errorf("failed to marshal aggregate: %w", err)
	}
	
	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to save aggregate to redis: %w", err)
	}
	
	return nil
}

func (r *AggregateRepository) GetAggregate(ctx context.Context, orderID string) (*models.DeliveryAggregate, error) {
	key := r.getKey(orderID)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get aggregate from redis: %w", err)
	}
	
	var aggregate models.DeliveryAggregate
	if err := json.Unmarshal([]byte(data), &aggregate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal aggregate: %w", err)
	}
	
	return &aggregate, nil
}

func (r *AggregateRepository) DeleteAggregate(ctx context.Context, orderID string) error {
	key := r.getKey(orderID)
	
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete aggregate from redis: %w", err)
	}
	
	return nil
}

func (r *AggregateRepository) GetReadyToDeliverOrders(ctx context.Context) ([]*models.DeliveryAggregate, error) {
	// Scan for all delivery keys
	pattern := "delivery:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to scan redis keys: %w", err)
	}
	
	var readyOrders []*models.DeliveryAggregate
	
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip errors for individual keys
		}
		
		var aggregate models.DeliveryAggregate
		if err := json.Unmarshal([]byte(data), &aggregate); err != nil {
			continue // Skip invalid data
		}
		
		if aggregate.CanStartDelivery() {
			readyOrders = append(readyOrders, &aggregate)
		}
	}
	
	return readyOrders, nil
}

func (r *AggregateRepository) getKey(orderID string) string {
	return fmt.Sprintf("delivery:%s", orderID)
}