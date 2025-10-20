package repository

import (
	"context"
	"order-service/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Transaction(ctx context.Context, fn func(repo *OrderRepository) (interface{}, error)) (interface{}, error) {
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

func (r *OrderRepository) withTx(tx *gorm.DB) *OrderRepository {
	return &OrderRepository{db: tx}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) FindLatestOrderByPatientID(ctx context.Context, patientID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems.Medicine").Where("patient_id = ? AND status != 'pending'", patientID).Order("created_at DESC").First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems.Medicine").Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems.Medicine").Where("patient_id = ?", patientID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindAll(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems.Medicine").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Model(order).Updates(order).Error
}

func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Order{}).Error
}

func (r *OrderRepository) FindByStatus(ctx context.Context, status models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems.Medicine").Where("status = ?", status).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems.Medicine").Where("id IN ?", ids).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindByDoctorID(ctx context.Context, doctorID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems.Medicine").Where("doctor_id = ? AND status = 'pending'", doctorID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
