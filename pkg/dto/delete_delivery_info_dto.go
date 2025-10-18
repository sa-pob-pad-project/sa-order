package dto

type DeleteDeliveryInfoRequestDto struct {
	ID string `json:"id" validate:"required"`
}

type DeleteDeliveryInfoResponseDto struct {
	ID        string `json:"id"`
	DeletedAt string `json:"deleted_at"`
}
