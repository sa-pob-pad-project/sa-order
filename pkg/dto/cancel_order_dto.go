package dto

type CancelOrderRequestDto struct {
	OrderID string `json:"order_id"`
}

type CancelOrderResponseDto struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
