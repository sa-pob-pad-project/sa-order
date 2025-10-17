package models

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	MedicineID uuid.UUID `gorm:"type:uuid;not null" json:"medicine_id"`
	Quantity   float64   `gorm:"type:numeric(12,2);not null;check:quantity > 0" json:"quantity"`
	Medicine   *Medicine `gorm:"foreignKey:MedicineID;references:ID" json:"medicine,omitempty"`
}

func (oi *OrderItem) TableName() string {
	return "order_items"
}
