package http

import (
	"fmt"
	"net/http"
	"shopping-service/internal/adapters/http/middleware"
	"shopping-service/internal/adapters/kafka"
	"shopping-service/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type AddToCartReq struct {
	ProductID string `json:"product_id" binding:"required"`
	Qty       int    `json:"qty" binding:"required"`
}

type CartHandler struct {
	cartSvc       *service.CartService
	productSvc    *service.ProductService
	kafkaProducer *kafka.KafkaProducer
}

func NewCartHandler(rg *gin.RouterGroup, cs *service.CartService, ps *service.ProductService, kp *kafka.KafkaProducer) {
	h := &CartHandler{cartSvc: cs, productSvc: ps, kafkaProducer: kp}
	rg.POST("/cart", h.AddToCart)
	rg.GET("/cart", h.GetCart)
	rg.PUT("/cart", h.UpdateCartItem)                // Update quantity
	rg.DELETE("/cart/:product_id", h.RemoveFromCart) // Remove item
	rg.POST("/checkout", h.Checkout)                 // protected: creates order event
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user"})
		return
	}
	var req AddToCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.cartSvc.AddToCart(uid, req.ProductID, req.Qty); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	// Kafka-Nachricht senden wenn Produkt zum Warenkorb hinzugefügt wird
	cartEvent := map[string]interface{}{
		"event_type": "item_added_to_cart",
		"user_id":    uid,
		"product_id": req.ProductID,
		"quantity":   req.Qty,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	if err := h.kafkaProducer.SendMessage(cartEvent); err != nil {
		// Log error but don't fail the request - cart was already updated
		fmt.Printf("Failed to send cart event to Kafka: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"status": "added", "product_id": req.ProductID, "qty": req.Qty})
}

func (h *CartHandler) GetCart(c *gin.Context) {
	uid, _ := middleware.GetUserID(c)
	cart, err := h.cartSvc.GetCart(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) Checkout(c *gin.Context) {
	uid, _ := middleware.GetUserID(c)
	cart, err := h.cartSvc.GetCart(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cart"})
		return
	}

	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	// Berechne Gesamtsumme und sammle Produktdetails
	var amount float64
	orderItems := []map[string]interface{}{}
	productIds := []string{}

	for _, item := range cart.Items {
		productIds = append(productIds, item.ProductID)
		
		// Echten Produktpreis aus DB holen
		product, err := h.productSvc.GetProductByID(item.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to get product %s: %v", item.ProductID, err),
			})
			return
		}
		
		itemPrice := product.Price * float64(item.Qty)
		amount += itemPrice

		orderItems = append(orderItems, map[string]interface{}{
			"product_id":   item.ProductID,
			"product_name": product.Name,
			"quantity":     item.Qty,
			"unit_price":   product.Price,
			"total_price":  itemPrice,
		})
	}

	// Erstelle umfassende Bestellnachricht
	order := map[string]interface{}{
		"event_type":   "order_created",
		"order_id":     fmt.Sprintf("ord-%d", time.Now().UnixNano()),
		"user_id":      uid,
		"items":        orderItems,
		"product_ids":  productIds,
		"total_amount": amount,
		"currency":     "EUR",
		"status":       "pending",
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	if err := h.kafkaProducer.SendMessage(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to produce order"})
		return
	}

	// Warenkorb nach erfolgreichem Checkout leeren
	_ = h.cartSvc.ClearCart(uid)

	c.JSON(http.StatusOK, gin.H{
		"status":       "order_created",
		"order_id":     order["order_id"],
		"total_amount": amount,
		"items_count":  len(cart.Items),
	})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user"})
		return
	}

	var req AddToCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Qty <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be greater than 0"})
		return
	}

	if err := h.cartSvc.UpdateCartItem(uid, req.ProductID, req.Qty); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cart item"})
		return
	}

	// Kafka-Event für Warenkorb-Update
	updateEvent := map[string]interface{}{
		"event_type":   "cart_item_updated",
		"user_id":      uid,
		"product_id":   req.ProductID,
		"new_quantity": req.Qty,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	_ = h.kafkaProducer.SendMessage(updateEvent)

	c.JSON(http.StatusOK, gin.H{"message": "cart item updated successfully"})
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user"})
		return
	}

	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id required"})
		return
	}

	if err := h.cartSvc.RemoveFromCart(uid, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove from cart"})
		return
	}

	// Kafka-Event für Item-Entfernung
	removeEvent := map[string]interface{}{
		"event_type": "item_removed_from_cart",
		"user_id":    uid,
		"product_id": productID,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	_ = h.kafkaProducer.SendMessage(removeEvent)

	c.JSON(http.StatusOK, gin.H{"message": "item removed from cart successfully"})
}
