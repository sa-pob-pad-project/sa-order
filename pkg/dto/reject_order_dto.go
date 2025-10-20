package dto

type RejectOrderRequestDto struct {
	OrderID string `json:"order_id"`
}

type RejectOrderResponseDto struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
