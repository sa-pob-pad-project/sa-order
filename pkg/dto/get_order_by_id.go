package dto

type OrderItem struct {
	MedicineID   string  `json:"medicine_id"`
	MedicineName string  `json:"medicine_name"`
	Quantity     float64 `json:"quantity"`
}

type GetOrderByIDResponseDto struct {
	OrderID        string      `json:"order_id"`
	PatientID      string      `json:"patient_id"`
	DoctorID       string      `json:"doctor_id"`
	TotalAmount    float64     `json:"total_amount"`
	Note           *string     `json:"note"`
	SubmittedAt    *string     `json:"submitted_at"`
	ReviewedAt     *string     `json:"reviewed_at"`
	Status         string      `json:"status"`
	DeliveryStatus *string     `json:"delivery_status"`
	DeliveryAt     *string     `json:"delivery_at"`
	OrderItems     []OrderItem `json:"order_items"`
}
