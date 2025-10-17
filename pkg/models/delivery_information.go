package models

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryMethodEnum string

const (
	DeliveryMethodFlash  DeliveryMethodEnum = "flash"
	DeliveryMethodPickUp DeliveryMethodEnum = "pick_up"
)

type DeliveryInformation struct {
	ID             uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID          `gorm:"type:uuid;not null" json:"user_id"`
	Address        string             `gorm:"type:text;not null" json:"address"`
	PhoneNumber    string             `gorm:"type:text;not null" json:"phone_number"`
	Version        int                `gorm:"type:int;not null;default:1;check:version > 0" json:"version"`
	DeliveryMethod DeliveryMethodEnum `gorm:"type:delivery_method_enum;not null" json:"delivery_method"`
	CreatedAt      time.Time          `gorm:"autoCreateTime:milli" json:"created_at"`
}

func (di *DeliveryInformation) TableName() string {
	return "delivery_informations"
}
