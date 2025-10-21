package dto

type PayOrderRequestDto struct {
	OrderID string `json:"order_id"`
}

type PayOrderResponseDto struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
