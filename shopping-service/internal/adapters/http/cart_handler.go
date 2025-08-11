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
	kafkaProducer *kafka.KafkaProducer
}

func NewCartHandler(rg *gin.RouterGroup, cs *service.CartService, kp *kafka.KafkaProducer) {
	h := &CartHandler{cartSvc: cs, kafkaProducer: kp}
	rg.POST("/cart", h.AddToCart)
	rg.GET("/cart", h.GetCart)
	rg.POST("/checkout", h.Checkout) // protected: creates order event
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
	c.JSON(http.StatusOK, gin.H{"status": "added"})
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
	// compute total (simple placeholder)
	var amount float64
	ids := []string{}
	for _, it := range cart.Items {
		ids = append(ids, it.ProductID)
		amount += float64(it.Qty) * 1.0
	} // real code: look up product price
	order := map[string]interface{}{
		"order_id":    fmt.Sprintf("ord-%d", time.Now().UnixNano()),
		"user_id":     uid,
		"product_ids": ids,
		"amount":      amount,
	}
	if err := h.kafkaProducer.SendMessage(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to produce order"})
		return
	}
	// clear cart on success
	_ = h.cartSvc.ClearCart(uid)
	c.JSON(http.StatusOK, gin.H{"status": "order_sent", "order": order})
}
