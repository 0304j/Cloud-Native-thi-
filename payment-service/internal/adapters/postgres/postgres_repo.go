package postgres

import (
	"context"
	"errors"
	"payment-service/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"payment-service/internal/ports"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) ports.PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Save(ctx context.Context, payment *models.Payment) error {
	entity := FromDomain(payment)
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return err
	}

	domainModel := entity.ToDomain()
	*payment = *domainModel
	return nil
}

func (r *PaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	var entity PaymentEntity
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return entity.ToDomain(), nil
}

func (r *PaymentRepository) FindAll(ctx context.Context) ([]models.Payment, error) {
	var entities []PaymentEntity
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	result := make([]models.Payment, len(entities))
	for index, entity := range entities {
		domainModel := entity.ToDomain()
		result[index] = *domainModel
	}
	return result, nil
}
