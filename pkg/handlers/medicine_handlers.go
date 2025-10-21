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
// @Summary Get all available medicines
// @Description Retrieves a list of all available medicines in the system. This endpoint does not require authentication.
// @Tags medicines
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetAllMedicinesResponseDto "Medicines retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving medicines"
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
// @Summary Get medicine details by ID
// @Description Retrieves detailed information about a specific medicine identified by its ID. This endpoint does not require authentication.
// @Tags medicines
// @Accept json
// @Produce json
// @Param id path string true "Medicine ID (UUID)"
// @Success 200 {object} dto.GetMedicineByIDResponseDto "Medicine retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid or missing medicine ID"
// @Failure 404 {object} response.ErrorResponse "Medicine not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving medicine"
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
