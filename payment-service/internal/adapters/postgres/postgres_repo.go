package postgres

import (
	"context"
	"errors"
	"payment-service/internal/domain/models"
	"time"

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

func (r *PaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&PaymentEntity{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("payment not found")
	}
	return nil
}

func (r *PaymentRepository) FindByStatus(ctx context.Context, status models.Status) ([]models.Payment, error) {
	var entities []PaymentEntity
	if err := r.db.WithContext(ctx).Where("status = ?",
		string(status)).Find(&entities).Error; err != nil {
		return nil, err
	}

	result := make([]models.Payment, len(entities))
	for i, entity := range entities {
		domainModel := entity.ToDomain()
		result[i] = *domainModel
	}
	return result, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	entity := FromDomain(payment)
	entity.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Save(entity)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("payment not found")
	}
	domainModel := entity.ToDomain()
	*payment = *domainModel
	return nil
}
