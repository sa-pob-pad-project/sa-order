package client_dto

type GetPatientsByIDsRequestDto struct {
	PatientIDs []string `json:"patient_ids" validate:"required,min=1,dive,required,uuid"`
}

type GetPatientProfileResponseDto struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
}
