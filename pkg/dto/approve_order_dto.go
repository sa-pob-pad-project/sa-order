package dto

type ApproveOrderRequestDto struct {
	OrderID string `json:"order_id"`
}

type ApproveOrderResponseDto struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
