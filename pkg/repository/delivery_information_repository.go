package repository

import (
	"context"
	"order-service/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryInformationRepository struct {
	db *gorm.DB
}

func NewDeliveryInformationRepository(db *gorm.DB) *DeliveryInformationRepository {
	return &DeliveryInformationRepository{
		db: db,
	}
}

func (r *DeliveryInformationRepository) Transaction(ctx context.Context, fn func(repo *DeliveryInformationRepository) (interface{}, error)) (interface{}, error) {
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

func (r *DeliveryInformationRepository) withTx(tx *gorm.DB) *DeliveryInformationRepository {
	return &DeliveryInformationRepository{db: tx}
}

func (r *DeliveryInformationRepository) Create(ctx context.Context, deliveryInfo *models.DeliveryInformation) error {
	return r.db.WithContext(ctx).Create(deliveryInfo).Error
}

func (r *DeliveryInformationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.DeliveryInformation, error) {
	var deliveryInfo models.DeliveryInformation
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&deliveryInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &deliveryInfo, nil
}

func (r *DeliveryInformationRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.DeliveryInformation, error) {
	var deliveryInfos []models.DeliveryInformation
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("version DESC").Find(&deliveryInfos).Error; err != nil {
		return nil, err
	}
	return deliveryInfos, nil
}

func (r *DeliveryInformationRepository) FindByUserIDAndDeliveryMethod(ctx context.Context, userID uuid.UUID, method models.DeliveryMethodEnum) ([]models.DeliveryInformation, error) {
	var deliveryInfos []models.DeliveryInformation
	if err := r.db.WithContext(ctx).Where("user_id = ? AND delivery_method = ?", userID, method).Order("version DESC").Find(&deliveryInfos).Error; err != nil {
		return nil, err
	}
	return deliveryInfos, nil
}

func (r *DeliveryInformationRepository) FindLatestByUserIDAndDeliveryMethod(ctx context.Context, userID uuid.UUID, method models.DeliveryMethodEnum) (*models.DeliveryInformation, error) {
	var deliveryInfo models.DeliveryInformation
	if err := r.db.WithContext(ctx).Where("user_id = ? AND delivery_method = ?", userID, method).Order("version DESC").First(&deliveryInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &deliveryInfo, nil
}

func (r *DeliveryInformationRepository) FindAll(ctx context.Context) ([]models.DeliveryInformation, error) {
	var deliveryInfos []models.DeliveryInformation
	if err := r.db.WithContext(ctx).Find(&deliveryInfos).Error; err != nil {
		return nil, err
	}
	return deliveryInfos, nil
}

func (r *DeliveryInformationRepository) Update(ctx context.Context, deliveryInfo *models.DeliveryInformation) error {
	return r.db.WithContext(ctx).Model(deliveryInfo).Updates(deliveryInfo).Error
}

func (r *DeliveryInformationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.DeliveryInformation{}).Error
}
