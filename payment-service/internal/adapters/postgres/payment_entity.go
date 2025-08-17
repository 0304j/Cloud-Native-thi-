package postgres

import (
	"time"

	"payment-service/internal/domain/models"

	"github.com/google/uuid"
)

type PaymentEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Provider  string    `gorm:"size:32;not null"`
	Amount    float64   `gorm:"not null"`
	Currency  string    `gorm:"size:8;not null"`
	Status    string    `gorm:"size:16;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (PaymentEntity) TableName() string {
	return "payments"
}

func (e *PaymentEntity) ToDomain() *models.Payment {
	return &models.Payment{
		ID:        e.ID,
		OrderID:   e.OrderID,
		UserID:    e.UserID,
		Provider:  e.Provider,
		Amount:    e.Amount,
		Currency:  e.Currency,
		Status:    models.Status(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func FromDomain(p *models.Payment) *PaymentEntity {
	return &PaymentEntity{
		ID:        p.ID,
		OrderID:   p.OrderID,
		UserID:    p.UserID,
		Provider:  p.Provider,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
