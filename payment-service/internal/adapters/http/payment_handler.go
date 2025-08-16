package http

import (
	"context"
	"net/http"
	"payment-service/internal/adapters/postgres"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	Repo *postgres.PaymentRepository
}

func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	payments, err := h.Repo.GetPayments(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}
