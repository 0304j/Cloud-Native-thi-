package ports

import (
	"context"
	"payment-service/internal/domain/models"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Save(ctx context.Context, p *models.Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	FindAll(ctx context.Context) ([]models.Payment, error)
	Update(ctx context.Context, p *models.Payment) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByStatus(ctx context.Context, status models.Status) ([]models.Payment, error)
}
