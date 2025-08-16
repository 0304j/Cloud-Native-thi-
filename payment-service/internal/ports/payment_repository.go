package ports

import (
	"context"
	"payment-service/internal/domain"
)

type PaymentRepository interface {
	GetPayment(ctx context.Context, id string) (domain.Payment, error)
	GetAllPayments(ctx context.Context) ([]domain.Payment, error)
	CreatePayment(ctx context.Context, payment domain.Payment) error
}

type PaymentService interface {
	GetPayment(ctx context.Context, id string) (domain.Payment, error)
	GetAllPayments(ctx context.Context) ([]domain.Payment, error)
	CreatePayment(ctx context.Context, payment domain.Payment) error
}
