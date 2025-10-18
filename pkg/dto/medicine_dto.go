package dto

type MedicineResponseDto struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Stock     float64 `json:"stock"`
	Unit      string  `json:"unit"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type GetAllMedicinesResponseDto struct {
	Medicines []MedicineResponseDto `json:"medicines"`
	Total     int                   `json:"total"`
}

type GetMedicineByIDResponseDto struct {
	Medicine MedicineResponseDto `json:"medicine"`
}
