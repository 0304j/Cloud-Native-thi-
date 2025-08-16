package http

import (
	"context"
	"fmt"
	"net/http"
	"payment-service/internal/domain"
	"payment-service/internal/ports"

	"github.com/gin-gonic/gin"
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
	id := c.Param("id")
	fmt.Println("Fetching payment with ID:", id)
	payment, err := h.Service.GetPayment(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var payment domain.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreatePayment(context.Background(), payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, payment)
}
