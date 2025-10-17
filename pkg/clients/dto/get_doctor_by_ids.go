package client_dto

type GetDoctorsByIDsRequestDto struct {
	DoctorIDs []string `json:"doctor_ids" validate:"required,min=1,dive,required,uuid"`
}

type GetDoctorProfileResponseDto struct {
	ID              string  `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Gender          string  `json:"gender"`
	PhoneNumber     string  `json:"phone_number"`
	Username        string  `json:"username"`
	Specialty       *string `json:"specialty,omitempty"`
	Bio             *string `json:"bio,omitempty"`
	YearsExperience *int    `json:"years_experience,omitempty"`
}
