package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shopping-service/internal/adapters/kafka"
	"shopping-service/internal/service"
)

// Beispiel-Order-Struktur
type Order struct {
	UserID     string   `json:"user_id"`
	ProductIDs []string `json:"product_ids"`
	Amount     float64  `json:"amount"`
	// weitere Felder nach Bedarf
}

// ListProductsHandler
func ListProductsHandler(service *service.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := service.GetAllProducts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

// RegisterHandler (dummy, implement user creation or proxy to auth-service)
func RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	}
}

// LoginHandler (dummy, implement JWT creation or proxy to auth-service)
func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"token": "jwt-token"})
	}
}

// ParseOrderFromRequest
func ParseOrderFromRequest(c *gin.Context) Order {
	var order Order
	_ = c.ShouldBindJSON(&order)
	return order
}

// SendOrder (Kafka Producer)
func SendOrder(producer *kafka.KafkaProducer, order Order) error {
	return producer.SendMessage(order)
}
