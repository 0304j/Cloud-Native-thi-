package service

import (
	"context"
	"fmt"
	"payment-service/internal/domain/models"
	"payment-service/internal/ports"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo ports.PaymentRepository
}

func NewService(r ports.PaymentRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreatePayment(ctx context.Context, p models.Payment) (*models.Payment, error) {
	if p.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	p.Status = models.StatusPending

	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	if err := s.repo.Save(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Service) GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) GetAllPayments(ctx context.Context) ([]models.Payment, error) {
	return s.repo.FindAll(ctx)
}
