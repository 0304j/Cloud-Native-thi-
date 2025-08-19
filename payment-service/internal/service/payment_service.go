package service

import (
	"context"
	"fmt"
	"log"
	"payment-service/internal/domain/models"
	"payment-service/internal/ports"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo           ports.PaymentRepository
	eventPublisher ports.EventPublisher
}

func NewService(r ports.PaymentRepository, ep ports.EventPublisher) *PaymentService {
	return &PaymentService{
		repo:           r,
		eventPublisher: ep,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, p models.Payment) (*models.Payment, error) {
	payment, err := models.NewPayment(p.OrderID, p.UserID, p.Provider,
		p.Amount, p.Currency)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, payment); err != nil {
		return nil, err
	}
	return payment, nil
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

func (s *PaymentService) ProcessOrderPayment(ctx context.Context, orderEvent models.OrderCreatedEvent) (*models.Payment, error) {
	log.Printf("Processing payment for order %s, amount: %.2f %s", orderEvent.OrderID,
		orderEvent.TotalAmount, orderEvent.Currency)

	orderID, err := uuid.Parse(orderEvent.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	userID, err := uuid.Parse(orderEvent.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	payment := models.Payment{
		OrderID:  orderID,
		UserID:   userID,
		Provider: orderEvent.PaymentProvider,
		Amount:   orderEvent.TotalAmount,
		Currency: orderEvent.Currency,
	}

	createdPayment, err := s.CreatePayment(ctx, payment)
	if err != nil {
		failedEvent := &models.PaymentFailedEvent{
			EventType: "payment_failed",
			OrderID:   orderEvent.OrderID,
			Reason:    err.Error(),
			Amount:    orderEvent.TotalAmount,
			Currency:  orderEvent.Currency,
			Timestamp: time.Now(),
		}
		s.eventPublisher.PublishPaymentFailed(ctx, failedEvent)
		return nil, err
	}

	// Simulate payment processing (in real world, this would call external payment provider)

	// Different delays by provider
	var processingTime time.Duration
	switch createdPayment.Provider {
	case "stripe":
		processingTime = 2 * time.Second
	case "paypal":
		processingTime = 4 * time.Second
	case "visa":
		processingTime = 3 * time.Second
	default:
		processingTime = 3 * time.Second
	}

	log.Printf("Processing payment with %s provider, will take %v...", createdPayment.Provider, processingTime)
	time.Sleep(processingTime)
	log.Printf("Payment processing delay completed for %s", createdPayment.Provider)

	// Update payment status to success
	updatedPayment, err := s.UpdatePaymentStatus(ctx, createdPayment.ID, models.StatusSuccess)
	if err != nil {
		log.Printf("Failed to update payment status to success: %v", err)
		// Still publish failed event
		failedEvent := &models.PaymentFailedEvent{
			EventType: "payment_failed",
			OrderID:   orderEvent.OrderID,
			Reason:    fmt.Sprintf("failed to update payment status: %v", err),
			Amount:    orderEvent.TotalAmount,
			Currency:  orderEvent.Currency,
			Timestamp: time.Now(),
		}
		s.eventPublisher.PublishPaymentFailed(ctx, failedEvent)
		return createdPayment, err
	}

	// Publish payment succeeded event
	successEvent := &models.PaymentSucceededEvent{
		EventType: "payment_succeeded",
		OrderID:   orderEvent.OrderID,
		PaymentID: updatedPayment.ID.String(),
		Amount:    updatedPayment.Amount,
		Currency:  updatedPayment.Currency,
		Timestamp: time.Now(),
	}

	if err := s.eventPublisher.PublishPaymentSucceeded(ctx, successEvent); err != nil {
		log.Printf("Failed to publish payment succeeded event: %v", err)
	}

	log.Printf("Payment %s successfully processed and marked as success", updatedPayment.ID.String())
	return updatedPayment, nil
}
