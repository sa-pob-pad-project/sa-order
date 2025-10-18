package dto

import (
	"order-service/pkg/models"

	"github.com/google/uuid"
)

type CreateDeliveryInfoRequestDto struct {
	Address        string                    `json:"address" validate:"required"`
	PhoneNumber    string                    `json:"phone_number" validate:"required"`
	DeliveryMethod models.DeliveryMethodEnum `json:"delivery_method" validate:"required,oneof=flash pick_up"`
}

type CreateDeliveryInfoResponseDto struct {
	DeliveryInfo DeliveryInfoDto `json:"delivery_info"`
}

// Conversion functions
func ToDeliveryInfoDto(info *models.DeliveryInformation) DeliveryInfoDto {
	return DeliveryInfoDto{
		ID:             info.ID.String(),
		UserID:         info.UserID.String(),
		Address:        info.Address,
		PhoneNumber:    info.PhoneNumber,
		Version:        info.Version,
		DeliveryMethod: info.DeliveryMethod,
		CreatedAt:      info.CreatedAt,
	}
}

func ToDeliveryInfoDtoList(infos []models.DeliveryInformation) []DeliveryInfoDto {
	result := make([]DeliveryInfoDto, len(infos))
	for i, info := range infos {
		result[i] = ToDeliveryInfoDto(&info)
	}
	return result
}

func ToDeliveryInformation(userID string, dto CreateDeliveryInfoRequestDto) (*models.DeliveryInformation, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return &models.DeliveryInformation{
		ID:             uuid.New(),
		UserID:         userUUID,
		Address:        dto.Address,
		PhoneNumber:    dto.PhoneNumber,
		Version:        1,
		DeliveryMethod: dto.DeliveryMethod,
	}, nil
}
