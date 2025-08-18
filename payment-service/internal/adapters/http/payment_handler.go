package http

import (
	"context"
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

	uid, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	payment, err := h.Service.GetPayment(context.Background(), uid)
	if err != nil {
		if err.Error() == "payment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
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
	c.JSON(http.StatusCreated, created)
}

func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var request struct {
		Status models.Status `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.Service.UpdatePaymentStatus(context.Background(), uid, request.Status)
	if err != nil {
		if err.Error() == "payment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := h.Service.DeletePayment(context.Background(), uid); err != nil {
		if err.Error() == "payment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *PaymentHandler) GetPaymentsByStatus(c *gin.Context) {
	statusStr := c.Query("status")
	if statusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status query parameter is required"})
		return
	}

	status := models.Status(statusStr)
	payments, err := h.Service.GetPaymentsByStatus(context.Background(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}
