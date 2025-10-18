package handlers

import (
	"order-service/pkg/apperr"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/response"
	service "order-service/pkg/services"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderService *service.OrderService
	// deliveryService *service.DeliveryService
}

func NewOrderHandler(orderService *service.OrderService, deliveryService *service.DeliveryService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		// deliveryService: deliveryService,
	}
}

// Handler functions
// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order in the system
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body dto.CreateOrderRequestDto true "Order creation data"
// @Success 201 {object} dto.CreateOrderResponseDto "Order created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to create order"
// @Router /api/order/v1/orders [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var body dto.CreateOrderRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}
	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.CreateOrder(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return response.Created(c, res)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Retrieve order details including order items and medicine information
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.GetOrderByIDResponseDto "Order retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order"
// @Router /api/order/v1/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	if orderID == "" {
		return response.BadRequest(c, "Order ID is required")
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetAllOrdersHistory godoc
// @Summary Get all orders history for the current patient
// @Description Retrieve all orders for the authenticated patient from JWT token
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.GetAllOrdersHistoryListDto "Orders retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve orders"
// @Router /api/order/v1/history [get]
// @Security Bearer
func (h *OrderHandler) GetAllOrdersHistory(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.GetAllOrdersHistoryByPatientID(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an existing order (doctor only - can only cancel their own orders)
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body dto.CancelOrderRequestDto true "Order ID to cancel"
// @Success 200 {object} dto.CancelOrderResponseDto "Order cancelled successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or order ID"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only doctors can cancel their own orders"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Failed to cancel order"
// @Router /api/order/v1/orders [delete]
// @Security Bearer
func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	var body dto.CancelOrderRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	if body.OrderID == "" {
		return response.BadRequest(c, "Order ID is required")
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.CancelOrder(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
