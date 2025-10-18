package dto

type GetAllOrdersHistoryResponseDto struct {
	OrderID        string      `json:"order_id"`
	PatientID      string      `json:"patient_id"`
	DoctorID       *string     `json:"doctor_id"`
	TotalAmount    float64     `json:"total_amount"`
	Note           *string     `json:"note"`
	SubmittedAt    *string     `json:"submitted_at"`
	ReviewedAt     *string     `json:"reviewed_at"`
	Status         string      `json:"status"`
	DeliveryStatus *string     `json:"delivery_status"`
	DeliveryAt     *string     `json:"delivery_at"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
	OrderItems     []OrderItem `json:"order_items"`
}

type GetAllOrdersHistoryListDto struct {
	Orders []GetAllOrdersHistoryResponseDto `json:"orders"`
	Total  int                              `json:"total"`
}
