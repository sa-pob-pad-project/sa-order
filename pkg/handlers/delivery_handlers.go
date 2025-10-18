package handlers

import (
	"order-service/pkg/apperr"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	service "order-service/pkg/services"

	"github.com/gofiber/fiber/v2"
)

type DeliveryInfoHandler struct {
	deliveryService *service.DeliveryService
}

func NewDeliveryInfoHandler(deliveryService *service.DeliveryService) *DeliveryInfoHandler {
	return &DeliveryInfoHandler{
		deliveryService: deliveryService,
	}
}

// CreateDeliveryInfo godoc
// @Summary Create a new delivery information
// @Description Create a new delivery information record
// @Tags delivery-info
// @Accept json
// @Produce json
// @Param delivery_info body dto.CreateDeliveryInfoRequestDto true "Delivery Information Data"
// @Success 201 {object} dto.CreateDeliveryInfoResponseDto "Delivery information created successfully"
// @Failure 400 {object} fiber.Map "Invalid request body"
// @Failure 500 {object} fiber.Map "Failed to create delivery information"
// @Router /api/delivery-info/v1 [post]
// @Security Bearer
func (h *DeliveryInfoHandler) CreateDeliveryInfo(c *fiber.Ctx) error {
	var body dto.CreateDeliveryInfoRequestDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.deliveryService.CreateDeliveryInfo(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetDeliveryInfo godoc
// @Summary Get delivery information by ID
// @Description Retrieve delivery information details by ID
// @Tags delivery-info
// @Accept json
// @Produce json
// @Param id path string true "Delivery Information ID"
// @Success 200 {object} dto.GetDeliveryInfoResponseDto "Delivery information retrieved successfully"
// @Failure 400 {object} fiber.Map "Invalid delivery information ID"
// @Failure 404 {object} fiber.Map "Delivery information not found"
// @Failure 500 {object} fiber.Map "Failed to retrieve delivery information"
// @Router /api/delivery-info/v1/{id} [get]
// @Security Bearer
func (h *DeliveryInfoHandler) GetDeliveryInfo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Delivery information ID is required"})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.deliveryService.GetDeliveryInfoByID(ctx, id)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetAllDeliveryInfos godoc
// @Summary Get all delivery information
// @Description Retrieve all delivery information records
// @Tags delivery-info
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetAllDeliveryInfosResponseDto "Delivery information retrieved successfully"
// @Failure 500 {object} fiber.Map "Failed to retrieve delivery information"
// @Router /api/delivery-info/v1 [get]
// @Security Bearer
func (h *DeliveryInfoHandler) GetAllDeliveryInfos(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	res, err := h.deliveryService.GetAllDeliveryInfos(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
