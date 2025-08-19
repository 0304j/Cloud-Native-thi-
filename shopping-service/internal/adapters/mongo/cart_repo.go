package mongo

import (
	"context"
	"shopping-service/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CartRepo struct{ coll *mongo.Collection }

func NewCartRepo(db *mongo.Database) *CartRepo { return &CartRepo{coll: db.Collection("carts")} }

func (r *CartRepo) AddItem(userID, productID string, qty int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"user_id": userID}
	update := bson.M{"$push": bson.M{"items": bson.M{"product_id": productID, "qty": qty}}}
	_, err := r.coll.UpdateOne(ctx, filter, update, &options.UpdateOptions{Upsert: ptrBool(true)})
	return err
}

func (r *CartRepo) GetCart(userID string) (domain.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var doc struct {
		UserID string `bson:"user_id"`
		Items  []struct {
			ProductID string `bson:"product_id"`
			Qty       int    `bson:"qty"`
		} `bson:"items"`
	}
	err := r.coll.FindOne(ctx, bson.M{"user_id": userID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Cart{UserID: userID, Items: []domain.CartItem{}}, nil
		}
		return domain.Cart{}, err
	}
	cart := domain.Cart{UserID: doc.UserID}
	for _, it := range doc.Items {
		cart.Items = append(cart.Items, domain.CartItem{ProductID: it.ProductID, Qty: it.Qty})
	}
	return cart, nil
}

func (r *CartRepo) ClearCart(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.coll.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}

func (r *CartRepo) UpdateItem(userID, productID string, qty int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":          userID,
		"items.product_id": productID,
	}
	update := bson.M{
		"$set": bson.M{
			"items.$.qty": qty,
		},
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		// Item nicht gefunden, füge es hinzu
		return r.AddItem(userID, productID, qty)
	}

	return nil
}

func (r *CartRepo) RemoveItem(userID, productID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{"product_id": productID},
		},
	}

	_, err := r.coll.UpdateOne(ctx, filter, update)
	return err
}

func ptrBool(b bool) *bool { return &b }
