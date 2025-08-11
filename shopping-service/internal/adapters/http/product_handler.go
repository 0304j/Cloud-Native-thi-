package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shopping-service/internal/adapters/http/middleware"
	"shopping-service/internal/adapters/kafka"
	"shopping-service/internal/domain"
	"shopping-service/internal/ports"
)

type ProductHandler struct {
	service       ports.ProductService
	kafkaProducer *kafka.KafkaProducer
}

func NewProductHandler(r *gin.Engine, service ports.ProductService, producer *kafka.KafkaProducer) {
	handler := &ProductHandler{
		service:       service,
		kafkaProducer: producer,
	}

	// Öffentliche Route
	r.GET("/products", handler.ListProducts)

	// Admin-Gruppe: Produkte anlegen (nur mit JWT und admin-Rolle)
	adminGroup := r.Group("/")
	adminGroup.Use(middleware.JWTMiddleware())
	adminGroup.Use(middleware.RequireRole("admin"))
	adminGroup.POST("/products", handler.CreateProduct)

	// User-Gruppe: Cart & Checkout (nur mit JWT, Rolle egal)
	userGroup := r.Group("/")
	userGroup.Use(middleware.JWTMiddleware())
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// UserID aus JWT holen und ins Produkt übernehmen
	userID, exists := c.Get(middleware.ContextUserIDKey)
	if exists {
		product.UserID = userID.(string) // domain.Product braucht das Feld!
	}

	err := h.service.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving product"})
		return
	}

	// ➕ Send to Kafka mit Fehlerbehandlung
	if err := h.kafkaProducer.SendMessage(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kafka send error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
