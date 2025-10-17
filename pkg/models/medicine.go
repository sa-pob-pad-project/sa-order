package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Medicine struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"type:text;not null" json:"name"`
	Price     float64        `gorm:"type:numeric(12,2);not null;check:price >= 0" json:"price"`
	Stock     float64        `gorm:"type:numeric(12,2);not null;check:stock >= 0" json:"stock"`
	Unit      string         `gorm:"type:text;not null" json:"unit"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (m *Medicine) TableName() string {
	return "medicines"
}
