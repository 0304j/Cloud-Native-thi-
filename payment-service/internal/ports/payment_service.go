package ports

import (
	"context"
	"payment-service/internal/domain/models"

	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, p models.Payment) (*models.Payment, error)
	GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	GetAllPayments(ctx context.Context) ([]models.Payment, error)
}
