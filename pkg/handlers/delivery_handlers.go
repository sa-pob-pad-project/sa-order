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
// @Summary Create new delivery information
// @Description Creates a new delivery information record for an order. Contains details about the delivery method and address.
// @Tags delivery-info
// @Accept json
// @Produce json
// @Param request body dto.CreateDeliveryInfoRequestDto true "Delivery information request data"
// @Success 201 {object} dto.CreateDeliveryInfoResponseDto "Delivery information created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or missing required fields"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 500 {object} response.ErrorResponse "Internal server error while creating delivery information"
// @Router /api/delivery-info/v1 [post]
// @Security ApiKeyAuth
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
// @Description Retrieves detailed information about a specific delivery record identified by its ID.
// @Tags delivery-info
// @Accept json
// @Produce json
// @Param id path string true "Delivery Information ID (UUID)"
// @Success 200 {object} dto.GetDeliveryInfoResponseDto "Delivery information retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid or missing delivery information ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 404 {object} response.ErrorResponse "Delivery information not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving delivery information"
// @Router /api/delivery-info/v1/{id} [get]
// @Security ApiKeyAuth
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
// @Summary Get all delivery information records
// @Description Retrieves all delivery information records from the system.
// @Tags delivery-info
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetDeliveryInfoResponseDto "Delivery information retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving delivery information"
// @Router /api/delivery-info/v1 [get]
// @Security ApiKeyAuth
func (h *DeliveryInfoHandler) GetAllDeliveryInfos(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	res, err := h.deliveryService.GetAllDeliveryInfos(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateDeliveryInfo godoc
// @Summary Update delivery information
// @Description Updates an existing delivery information record with new data.
// @Tags delivery-info
// @Accept json
// @Produce json
// @Param request body dto.UpdateDeliveryInfoRequestDto true "Updated delivery information request data"
// @Success 200 {object} dto.UpdateDeliveryInfoResponseDto "Delivery information updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or missing delivery information ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 404 {object} response.ErrorResponse "Delivery information not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while updating delivery information"
// @Router /api/delivery-info/v1 [put]
// @Security ApiKeyAuth
func (h *DeliveryInfoHandler) UpdateDeliveryInfo(c *fiber.Ctx) error {
	var body dto.UpdateDeliveryInfoRequestDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.deliveryService.UpdateDeliveryInfo(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// DeleteDeliveryInfo godoc
// @Summary Delete delivery information
// @Description Deletes an existing delivery information record.
// @Tags delivery-info
// @Accept json
// @Produce json
// @Param request body dto.DeleteDeliveryInfoRequestDto true "Delivery information ID to delete"
// @Success 200 {object} dto.DeleteDeliveryInfoResponseDto "Delivery information deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid or missing delivery information ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 404 {object} response.ErrorResponse "Delivery information not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while deleting delivery information"
// @Router /api/delivery-info/v1 [delete]
// @Security ApiKeyAuth
func (h *DeliveryInfoHandler) DeleteDeliveryInfo(c *fiber.Ctx) error {
	var body dto.DeleteDeliveryInfoRequestDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.deliveryService.DeleteDeliveryInfo(ctx, body.ID)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetDeliveryInfosByMethod godoc
// @Summary Get delivery information by delivery method
// @Description Retrieves all delivery information records filtered by the specified delivery method (e.g., 'flash' for express delivery or 'pick_up' for customer pickup).
// @Tags delivery-info
// @Accept json
// @Produce json
// @Param method query string true "Delivery method filter: 'flash' (express delivery) or 'pick_up' (customer pickup)"
// @Success 200 {object} dto.GetAllDeliveryInfosResponseDto "Delivery information retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Missing or invalid delivery method query parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving delivery information"
// @Router /api/delivery-info/v1/methods [get]
// @Security ApiKeyAuth
func (h *DeliveryInfoHandler) GetDeliveryInfosByMethod(c *fiber.Ctx) error {
	method := c.Query("method")
	if method == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query parameter 'method' is required"})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.deliveryService.GetDeliveryInfosByMethod(ctx, method)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
