package repository

import (
	"context"
	"order-service/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicineRepository struct {
	db *gorm.DB
}

func NewMedicineRepository(db *gorm.DB) *MedicineRepository {
	return &MedicineRepository{
		db: db,
	}
}

func (r *MedicineRepository) Transaction(ctx context.Context, fn func(repo *MedicineRepository) (interface{}, error)) (interface{}, error) {
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

func (r *MedicineRepository) withTx(tx *gorm.DB) *MedicineRepository {
	return &MedicineRepository{db: tx}
}

func (r *MedicineRepository) Create(ctx context.Context, medicine *models.Medicine) error {
	return r.db.WithContext(ctx).Create(medicine).Error
}

func (r *MedicineRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Medicine, error) {
	var medicine models.Medicine
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&medicine).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &medicine, nil
}

func (r *MedicineRepository) FindAll(ctx context.Context) ([]models.Medicine, error) {
	var medicines []models.Medicine
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&medicines).Error; err != nil {
		return nil, err
	}
	return medicines, nil
}

func (r *MedicineRepository) Update(ctx context.Context, medicine *models.Medicine) error {
	return r.db.WithContext(ctx).Model(medicine).Updates(medicine).Error
}

func (r *MedicineRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Medicine{}).Error
}

func (r *MedicineRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Medicine, error) {
	var medicines []models.Medicine
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Where("deleted_at IS NULL").Find(&medicines).Error; err != nil {
		return nil, err
	}
	return medicines, nil
}
