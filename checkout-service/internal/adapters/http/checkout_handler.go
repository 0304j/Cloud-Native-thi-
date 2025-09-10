package http

import (
	"net/http"
	"time"

	"checkout-service/internal/domain/models"
	"checkout-service/internal/ports"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	orderService ports.OrderService
}

func NewCheckoutHandler(orderService ports.OrderService) *CheckoutHandler {
	return &CheckoutHandler{
		orderService: orderService,
	}
}

func (h *CheckoutHandler) CreateOrder(c *gin.Context) {
	var req models.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the checkout request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values
	req.EventType = "order_created"
	req.Status = "pending"
	req.Timestamp = time.Now().Format(time.RFC3339)

	// Process the order
	order := req.ToOrder()
	if err := h.orderService.ProcessOrder(c.Request.Context(), order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process order"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"status":       "order_created",
		"order_id":     order.ID.String(),
		"order_type":   order.OrderType,
		"total_amount": order.TotalAmount,
		"currency":     order.Currency,
		"created_at":   order.CreatedAt,
	})
}

func (h *CheckoutHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "checkout-service",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}