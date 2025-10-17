package models

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryStatus string

const (
	DeliveryStatusPending   DeliveryStatus = "pending"
	DeliveryStatusInTransit DeliveryStatus = "in_transit"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusFailed    DeliveryStatus = "failed"
)

type Delivery struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID             uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	DeliveryInformation uuid.UUID      `gorm:"type:uuid;not null" json:"delivery_information"`
	TrackingNumber      *string        `gorm:"type:text" json:"tracking_number,omitempty"`
	Status              DeliveryStatus `gorm:"type:delivery_status;not null;default:'pending'" json:"status"`
	DeliveredAt         *time.Time     `json:"delivered_at,omitempty"`
	CreatedAt           time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
}

func (d *Delivery) TableName() string {
	return "deliveries"
}
