package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusApproved   OrderStatus = "approved"
	OrderStatusRejected   OrderStatus = "rejected"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID               uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID        uuid.UUID   `gorm:"type:uuid;not null" json:"patient_id"`
	DoctorID         *uuid.UUID  `gorm:"type:uuid" json:"doctor_id,omitempty"`
	TotalAmount      float64     `gorm:"type:numeric(12,2);not null;check:total_amount >= 0" json:"total_amount"`
	Note             *string     `gorm:"type:text" json:"note,omitempty"`
	SubmittedAt      *time.Time  `json:"submitted_at,omitempty"`
	ReviewedAt       *time.Time  `json:"reviewed_at,omitempty"`
	Status           OrderStatus `gorm:"type:order_status;not null;default:'pending'" json:"status"`
	CreatedAt        time.Time   `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt        time.Time   `gorm:"autoUpdateTime:milli" json:"updated_at"`
	OrderItems       []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items,omitempty"`
}

func (o *Order) TableName() string {
	return "orders"
}
