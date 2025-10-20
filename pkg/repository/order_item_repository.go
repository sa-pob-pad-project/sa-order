package repository

import (
	"context"
	"order-service/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

func (r *OrderItemRepository) Transaction(ctx context.Context, fn func(repo *OrderItemRepository) (interface{}, error)) (interface{}, error) {
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

func (r *OrderItemRepository) withTx(tx *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{db: tx}
}

func (r *OrderItemRepository) Create(ctx context.Context, orderItem *models.OrderItem) error {
	return r.db.WithContext(ctx).Create(orderItem).Error
}

func (r *OrderItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&orderItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &orderItem, nil
}

func (r *OrderItemRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.OrderItem, error) {
	var items []models.OrderItem
	if err := r.db.WithContext(ctx).Preload("Medicine").Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderItemRepository) FindAll(ctx context.Context) ([]models.OrderItem, error) {
	var items []models.OrderItem
	if err := r.db.WithContext(ctx).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderItemRepository) Update(ctx context.Context, orderItem *models.OrderItem) error {
	return r.db.WithContext(ctx).Model(orderItem).Updates(orderItem).Error
}

func (r *OrderItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.OrderItem{}).Error
}

func (r *OrderItemRepository) DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("order_id = ?", orderID).Delete(&models.OrderItem{}).Error
}
