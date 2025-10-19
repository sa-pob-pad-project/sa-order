package dto

type PatientInfo struct {
	PatientID   string `json:"patient_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
}

type GetAllOrdersForDoctorResponseDto struct {
	OrderID        string       `json:"order_id"`
	PatientID      string       `json:"patient_id"`
	PatientInfo    *PatientInfo `json:"patient_info"`
	DoctorID       *string      `json:"doctor_id"`
	TotalAmount    float64      `json:"total_amount"`
	Note           *string      `json:"note"`
	SubmittedAt    *string      `json:"submitted_at"`
	ReviewedAt     *string      `json:"reviewed_at"`
	Status         string       `json:"status"`
	DeliveryStatus *string      `json:"delivery_status"`
	DeliveryAt     *string      `json:"delivery_at"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
	OrderItems     []OrderItem  `json:"order_items"`
}

type GetAllOrdersForDoctorListDto struct {
	Orders []GetAllOrdersForDoctorResponseDto `json:"orders"`
	Total  int                                `json:"total"`
}
