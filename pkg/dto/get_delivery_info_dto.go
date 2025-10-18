package dto

import (
	"order-service/pkg/models"
	"time"
)

type DeliveryInfoDto struct {
	ID             string                    `json:"id"`
	UserID         string                    `json:"user_id"`
	Address        string                    `json:"address"`
	PhoneNumber    string                    `json:"phone_number"`
	Version        int                       `json:"version"`
	DeliveryMethod models.DeliveryMethodEnum `json:"delivery_method"`
	CreatedAt      time.Time                 `json:"created_at"`
}

type GetDeliveryInfoResponseDto struct {
	DeliveryInfo DeliveryInfoDto `json:"delivery_info"`
}
