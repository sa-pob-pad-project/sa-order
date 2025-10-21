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
// @Description Creates a new order in the system. Only patients can create orders. The order will be assigned to the authenticated patient's doctor.
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequestDto true "Order creation request data"
// @Success 201 {object} dto.CreateOrderResponseDto "Order created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or missing required fields"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only patients can create orders"
// @Failure 500 {object} response.ErrorResponse "Internal server error while creating order"
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
// @Description Updates an order with new items or modifications. Only doctors can update orders they created. Supports adding, editing, or removing order items.
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.UpdateOrderRequestDto true "Order update request data"
// @Success 200 {object} dto.UpdateOrderResponseDto "Order updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or malformed order ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only the doctor who created this order can update it"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while updating order"
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
// @Description Retrieves detailed information about a specific order including all order items and associated medicine information.
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} dto.GetOrderByIDResponseDto "Order retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid or missing order ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving order"
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
// @Summary Get all orders for the current patient
// @Description Retrieves the complete order history for the authenticated patient. The patient is identified from the JWT authentication token.
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetAllOrdersHistoryListDto "Orders retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving orders"
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
// @Summary Get the latest order for the current patient
// @Description Retrieves the most recent order for the authenticated patient. Returns the latest order regardless of its status.
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetOrderByIDResponseDto "Order retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 404 {object} response.ErrorResponse "No orders found for this patient"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving order"
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
// @Summary Get the latest order for a specific patient
// @Description Retrieves the most recent order for a specified patient. Only the assigned doctor can access this endpoint. The doctor is verified through the JWT token.
// @Tags orders
// @Accept json
// @Produce json
// @Param patient_id path string true "Patient ID (UUID)"
// @Success 200 {object} dto.GetOrderByIDResponseDto "Order retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 403 {object} response.ErrorResponse "Forbidden - doctors can only access orders for their assigned patients"
// @Failure 404 {object} response.ErrorResponse "No orders found for this patient"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving order"
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
// @Summary Cancel an existing order
// @Description Cancels an order (doctor only). Only the doctor who created the order can cancel it. The order status will be changed to cancelled.
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.CancelOrderRequestDto true "Cancel order request data"
// @Success 200 {object} dto.CancelOrderResponseDto "Order cancelled successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or missing order ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only the doctor who created this order can cancel it"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while cancelling order"
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
// @Summary Approve an existing order
// @Description Approves an order and sets its status to approved (doctor only). Only the doctor who created the order can approve it.
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.ApproveOrderRequestDto true "Approve order request data"
// @Success 200 {object} dto.ApproveOrderResponseDto "Order approved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or missing order ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only the doctor who created this order can approve it"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while approving order"
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
// @Summary Reject an existing order
// @Description Rejects an order and sets its status to rejected (doctor only). Only the doctor who created the order can reject it.
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.RejectOrderRequestDto true "Reject order request data"
// @Success 200 {object} dto.RejectOrderResponseDto "Order rejected successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or missing order ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only the doctor who created this order can reject it"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while rejecting order"
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

func (h OrderHandler) PayOrder(c *fiber.Ctx) error {
	var body dto.PayOrderRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	if body.OrderID == "" {
		return response.BadRequest(c, "Order ID is required")
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.PayOrder(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// PayOrder godoc
// @Summary Mark an order as paid
// @Description Marks an order as paid and updates its payment status. Only the patient who created the order can pay it.
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.PayOrderRequestDto true "Pay order request data"
// @Success 200 {object} response.SuccessResponse "Order payment recorded successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or missing order ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only the patient can pay their own orders"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error while processing payment"
// @Router /api/order/v1/orders/pay [post]
// @Security ApiKeyAuth

// GetAllOrdersForDoctor godoc
// @Summary Get all orders for the current doctor
// @Description Retrieves all orders created by the authenticated doctor. Includes patient information for each order. The doctor is identified from the JWT authentication token.
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetAllOrdersForDoctorListDto "Orders retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only doctors can access this endpoint"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving orders"
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

// GetAllOrdersHistoryForDoctor godoc
// @Summary Get approved or rejected orders for the current doctor
// @Description Retrieves approved or rejected orders created by the authenticated doctor. Includes patient information for each order. Can filter by status using the optional query parameter. Valid status values are "approved" or "rejected".
// @Tags orders
// @Accept json
// @Produce json
// @Param status query string false "Filter by status: 'approved' or 'rejected'. If omitted, returns all approved and rejected orders."
// @Success 200 {object} dto.GetAllOrdersForDoctorListDto "Orders retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - authentication token missing or invalid"
// @Failure 403 {object} response.ErrorResponse "Forbidden - only doctors can access this endpoint"
// @Failure 500 {object} response.ErrorResponse "Internal server error while retrieving orders"
// @Router /api/order/v1/orders/doctor/history [get]
// @Security ApiKeyAuth
func (h *OrderHandler) GetAllOrdersHistoryForDoctor(c *fiber.Ctx) error {
	statusFilter := c.Query("status", "")

	ctx := contextUtils.GetContext(c)
	res, err := h.orderService.GetAllOrdersHistoryForDoctor(ctx, statusFilter)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
