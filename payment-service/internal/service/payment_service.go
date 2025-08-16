package service

import (
	"context"
	"payment-service/internal/domain"
	"payment-service/internal/ports"
)

type PaymentService struct {
	Repo ports.PaymentRepository
}

func (s *PaymentService) GetPayment(ctx context.Context, id string) (domain.Payment, error) {
	return s.Repo.GetPayment(ctx, id)
}

func (s *PaymentService) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	return s.Repo.GetAllPayments(ctx)
}

func (s *PaymentService) CreatePayment(ctx context.Context, payment domain.Payment) error {
	return s.Repo.CreatePayment(ctx, payment)
}
