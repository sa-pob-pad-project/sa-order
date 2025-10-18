package handlers

import (
	"order-service/pkg/apperr"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/response"
	service "order-service/pkg/services"

	"github.com/gofiber/fiber/v2"
)

type MedicineHandler struct {
	medicineService *service.MedicineService
}

func NewMedicineHandler(medicineService *service.MedicineService) *MedicineHandler {
	return &MedicineHandler{
		medicineService: medicineService,
	}
}

// GetAllMedicines godoc
// @Summary Get all medicines
// @Description Retrieve all available medicines from the system
// @Tags medicines
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.GetAllMedicinesResponseDto "Medicines retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve medicines"
// @Router /api/medicine/v1/medicines [get]
func (h *MedicineHandler) GetAllMedicines(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	res, err := h.medicineService.GetAllMedicines(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetMedicineByID godoc
// @Summary Get medicine by ID
// @Description Retrieve medicine details by ID
// @Tags medicines
// @Accept  json
// @Produce  json
// @Param id path string true "Medicine ID"
// @Success 200 {object} dto.GetMedicineByIDResponseDto "Medicine retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid medicine ID"
// @Failure 404 {object} response.ErrorResponse "Medicine not found"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve medicine"
// @Router /api/medicine/v1/medicines/{id} [get]
func (h *MedicineHandler) GetMedicineByID(c *fiber.Ctx) error {
	medicineID := c.Params("id")
	if medicineID == "" {
		return response.BadRequest(c, "Medicine ID is required")
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.medicineService.GetMedicineByID(ctx, medicineID)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
