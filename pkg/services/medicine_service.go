package service

import (
	"context"
	"order-service/pkg/apperr"
	"order-service/pkg/dto"
	"order-service/pkg/repository"

	"github.com/google/uuid"
)

type MedicineService struct {
	medicineRepository *repository.MedicineRepository
}

func NewMedicineService(medicineRepo *repository.MedicineRepository) *MedicineService {
	return &MedicineService{
		medicineRepository: medicineRepo,
	}
}

func (s *MedicineService) GetAllMedicines(ctx context.Context) (*dto.GetAllMedicinesResponseDto, error) {
	medicines, err := s.medicineRepository.FindAll(ctx)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "Failed to retrieve medicines", err)
	}

	medicineList := make([]dto.MedicineResponseDto, len(medicines))
	for i, medicine := range medicines {
		medicineList[i] = dto.MedicineResponseDto{
			ID:        medicine.ID.String(),
			Name:      medicine.Name,
			Price:     medicine.Price,
			Stock:     medicine.Stock,
			Unit:      medicine.Unit,
			CreatedAt: medicine.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: medicine.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &dto.GetAllMedicinesResponseDto{
		Medicines: medicineList,
		Total:     len(medicineList),
	}, nil
}

func (s *MedicineService) GetMedicineByID(ctx context.Context, medicineID string) (*dto.GetMedicineByIDResponseDto, error) {
	if medicineID == "" {
		return nil, apperr.New(apperr.CodeBadRequest, "Medicine ID is required", nil)
	}

	id, err := uuid.Parse(medicineID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "Invalid medicine ID format", err)
	}

	medicine, err := s.medicineRepository.FindByID(ctx, id)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "Medicine not found", err)
	}

	return &dto.GetMedicineByIDResponseDto{
		Medicine: dto.MedicineResponseDto{
			ID:        medicine.ID.String(),
			Name:      medicine.Name,
			Price:     medicine.Price,
			Stock:     medicine.Stock,
			Unit:      medicine.Unit,
			CreatedAt: medicine.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: medicine.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}
