package client_dto

type GetLatestAppointmentResponseDto struct {
	DoctorID        string `json:"doctor_id"`
	DoctorFirstName string `json:"doctor_first_name"`
	DoctorLastName  string `json:"doctor_last_name"`
	Specialty       string `json:"specialty"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	Status         string `json:"status"`
}
