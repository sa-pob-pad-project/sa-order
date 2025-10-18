package dto

import (
	"order-service/pkg/models"
)

type UpdateDeliveryInfoRequestDto struct {
	ID             string                    `json:"id" validate:"required"`
	Address        string                    `json:"address" validate:"required"`
	PhoneNumber    string                    `json:"phone_number" validate:"required"`
	DeliveryMethod models.DeliveryMethodEnum `json:"delivery_method" validate:"required,oneof=flash pick_up"`
}

type UpdateDeliveryInfoResponseDto struct {
	DeliveryInfo DeliveryInfoDto `json:"delivery_info"`
}
