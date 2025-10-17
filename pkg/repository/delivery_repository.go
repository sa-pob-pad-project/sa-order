package repository

import (
	"context"
	"order-service/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryRepository struct {
	db *gorm.DB
}

func NewDeliveryRepository(db *gorm.DB) *DeliveryRepository {
	return &DeliveryRepository{
		db: db,
	}
}

func (r *DeliveryRepository) Transaction(ctx context.Context, fn func(repo *DeliveryRepository) (interface{}, error)) (interface{}, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	repoWithTx := r.withTx(tx)

	result, err := fn(repoWithTx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *DeliveryRepository) withTx(tx *gorm.DB) *DeliveryRepository {
	return &DeliveryRepository{db: tx}
}

func (r *DeliveryRepository) Create(ctx context.Context, delivery *models.Delivery) error {
	return r.db.WithContext(ctx).Create(delivery).Error
}

func (r *DeliveryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Delivery, error) {
	var delivery models.Delivery
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&delivery).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &delivery, nil
}

func (r *DeliveryRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Delivery, error) {
	var delivery models.Delivery
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&delivery).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &delivery, nil
}

func (r *DeliveryRepository) FindAll(ctx context.Context) ([]models.Delivery, error) {
	var deliveries []models.Delivery
	if err := r.db.WithContext(ctx).Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (r *DeliveryRepository) Update(ctx context.Context, delivery *models.Delivery) error {
	return r.db.WithContext(ctx).Model(delivery).Updates(delivery).Error
}

func (r *DeliveryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Delivery{}).Error
}

func (r *DeliveryRepository) FindByStatus(ctx context.Context, status models.DeliveryStatus) ([]models.Delivery, error) {
	var deliveries []models.Delivery
	if err := r.db.WithContext(ctx).Where("status = ?", status).Order("created_at DESC").Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}
