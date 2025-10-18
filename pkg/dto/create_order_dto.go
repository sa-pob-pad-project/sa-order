package dto

type CreateOrderRequestDto struct {
	Note *string `json:"note"`
}

type CreateOrderResponseDto struct {
	OrderID string `json:"order_id"`
}
