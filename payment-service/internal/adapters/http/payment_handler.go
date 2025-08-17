package http

import (
	"context"
	"fmt"
	"net/http"
	"payment-service/internal/domain/models"
	"payment-service/internal/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	Service ports.PaymentService
}

func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	payments, err := h.Service.GetAllPayments(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	idStr := c.Param("id")
	fmt.Println("Fetching payment with ID:", idStr)

	uid, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	payment, err := h.Service.GetPayment(context.Background(), uid)
	if err != nil {
		// Optional: 404 sauber unterscheiden, falls dein Service/Repo ErrNotFound zurückgibt
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.Service.CreatePayment(context.Background(), payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Das vom Service/Repo aktualisierte Objekt (inkl. ID/Timestamps) zurückgeben
	c.JSON(http.StatusCreated, created)
}
