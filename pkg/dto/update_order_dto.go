package dto

import "github.com/google/uuid"

type OrderItemInput struct {
	MedicineID uuid.UUID `json:"medicine_id"`
	Quantity   float64   `json:"quantity"`
}

type UpdateOrderRequestDto struct {
	OrderID    string           `json:"order_id"`
	OrderItems []OrderItemInput `json:"order_items"`
}

type UpdateOrderResponseDto struct {
	OrderID string `json:"order_id"`
}
