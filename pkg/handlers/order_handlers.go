package handlers

import (
	"fmt"
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

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
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
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only patients can create orders"
// @Failure 500 {object} response.ErrorResponse "Failed to create order"
// @Router /api/order/v1/orders [post]
// @Security ApiKeyAuth
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

// UpdateOrder godoc
// @Summary Update an existing order
// @Description Update an order (doctor only - can only edit their own orders). Can add, edit, or remove order items.
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body dto.UpdateOrderRequestDto true "Order update data"
// @Success 200 {object} dto.UpdateOrderResponseDto "Order updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or order ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only doctors can update their own orders"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Failed to update order"
// @Router /api/order/v1/orders [put]
// @Security ApiKeyAuth
func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {

	var body dto.UpdateOrderRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.UpdateOrder(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
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
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order"
// @Router /api/order/v1/orders/{id} [get]
// @Security ApiKeyAuth
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	fmt.Println("id hit")
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
// @Router /api/order/v1/orders [get]
// @Security ApiKeyAuth
func (h *OrderHandler) GetAllOrdersHistory(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.GetAllOrdersHistoryByPatientID(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetLatestOrder godoc
// @Summary Get latest order for the current patient
// @Description Retrieve the most recent order for the authenticated patient
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.GetOrderByIDResponseDto "Order retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "No orders found"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order"
// @Router /api/order/v1/orders/latest [get]
// @Security ApiKeyAuth
func (h *OrderHandler) GetLatestOrder(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.GetLatestOrderByPatientID(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetLatestOrderByPatientID godoc
// @Summary Get latest order for a specific patient (doctor only)
// @Description Retrieve the most recent order for a specific patient. Only the assigned doctor can access this endpoint.
// @Tags orders
// @Accept  json
// @Produce  json
// @Param patient_id path string true "Patient ID"
// @Success 200 {object} dto.GetOrderByIDResponseDto "Order retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - doctor can only access their own patient's orders"
// @Failure 404 {object} response.ErrorResponse "No orders found for this patient"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order"
// @Router /api/order/v1/orders/latest/{patient_id} [get]
// @Security ApiKeyAuth
func (h *OrderHandler) GetLatestOrderByPatientID(c *fiber.Ctx) error {
	patientID := c.Params("patient_id")
	if patientID == "" {
		return response.BadRequest(c, "Patient ID is required")
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.GetLatestOrderByPatientIDForDoctor(ctx, patientID)
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
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only doctors can cancel their own orders"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Failed to cancel order"
// @Router /api/order/v1/orders [delete]
// @Security ApiKeyAuth
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

// ApproveOrder godoc
// @Summary Approve an order
// @Description Approve an existing order (doctor only - can only approve their own orders). Sets order status to approved.
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body dto.ApproveOrderRequestDto true "Order ID to approve"
// @Success 200 {object} dto.ApproveOrderResponseDto "Order approved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or order ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only doctors can approve their own orders"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Failed to approve order"
// @Router /api/order/v1/orders/confirm [post]
// @Security ApiKeyAuth
func (h *OrderHandler) ApproveOrder(c *fiber.Ctx) error {
	var body dto.ApproveOrderRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	if body.OrderID == "" {
		return response.BadRequest(c, "Order ID is required")
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.ApproveOrder(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// RejectOrder godoc
// @Summary Reject an order
// @Description Reject an existing order (doctor only - can only reject their own orders). Sets order status to rejected.
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body dto.RejectOrderRequestDto true "Order ID to reject"
// @Success 200 {object} dto.RejectOrderResponseDto "Order rejected successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or order ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only doctors can reject their own orders"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Failed to reject order"
// @Router /api/order/v1/orders/reject [post]
// @Security ApiKeyAuth
func (h *OrderHandler) RejectOrder(c *fiber.Ctx) error {
	var body dto.RejectOrderRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	if body.OrderID == "" {
		return response.BadRequest(c, "Order ID is required")
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.RejectOrder(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetAllOrdersForDoctor godoc
// @Summary Get all orders for the current doctor
// @Description Retrieve all orders created by the authenticated doctor with patient information
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.GetAllOrdersForDoctorListDto "Orders retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only doctors can access this endpoint"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve orders"
// @Router /api/order/v1/orders/doctor [get]
// @Security ApiKeyAuth
func (h *OrderHandler) GetAllOrdersForDoctor(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.GetAllOrdersByDoctorID(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
