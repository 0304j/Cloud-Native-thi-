package service

import (
	"context"
	"fmt"
	"payment-service/internal/domain/models"
	"payment-service/internal/ports"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo ports.PaymentRepository
}

func NewService(r ports.PaymentRepository) *PaymentService {
	return &PaymentService{repo: r}
}

func (s *PaymentService) CreatePayment(ctx context.Context, p models.Payment) (*models.Payment, error) {
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

func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *PaymentService) GetAllPayments(ctx context.Context) ([]models.Payment, error) {
	return s.repo.FindAll(ctx)
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status models.Status) (*models.Payment, error) {
	payment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	switch status {
	case models.StatusSuccess:
		if err := payment.MarkSuccess(); err != nil {
			return nil, err
		}
	case models.StatusFailed:
		if err := payment.MarkFailed(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid status transition to: %s", status)
	}

	if err := s.repo.Update(ctx, payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *PaymentService) DeletePayment(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *PaymentService) GetPaymentsByStatus(ctx context.Context, status models.Status) ([]models.Payment, error) {
	return s.repo.FindByStatus(ctx, status)
}
