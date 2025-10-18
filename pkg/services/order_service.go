package service

import (
	"context"
	"order-service/pkg/apperr"
	"order-service/pkg/clients"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/models"
	"order-service/pkg/repository"
	"order-service/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	db                     *gorm.DB
	orderRepository        *repository.OrderRepository
	orderItemRepository    *repository.OrderItemRepository
	medicineRepository     *repository.MedicineRepository
	deliveryRepository     *repository.DeliveryRepository
	deliveryInfoRepository *repository.DeliveryInformationRepository
	userClient             *clients.UserClient
	appointmentClient      *clients.AppointmentClient
}

func NewOrderService(
	db *gorm.DB,
	orderRepo *repository.OrderRepository,
	orderItemRepo *repository.OrderItemRepository,
	medicineRepo *repository.MedicineRepository,
	deliveryRepo *repository.DeliveryRepository,
	deliveryInfoRepo *repository.DeliveryInformationRepository,
	userClient *clients.UserClient,
	appointmentClient *clients.AppointmentClient,
) *OrderService {
	return &OrderService{
		db:                     db,
		orderRepository:        orderRepo,
		orderItemRepository:    orderItemRepo,
		medicineRepository:     medicineRepo,
		deliveryRepository:     deliveryRepo,
		deliveryInfoRepository: deliveryInfoRepo,
		userClient:             userClient,
		appointmentClient:      appointmentClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, body dto.CreateOrderRequestDto) (*dto.CreateOrderResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "patient" {
		return nil, apperr.New(apperr.CodeForbidden, "only patients can create orders", nil)
	}
	patientID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}
	appointment, err := s.appointmentClient.GetLatestAppointmentByPatientID(ctx, patientID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to get latest appointment", err)
	}
	if appointment == nil {
		return nil, apperr.New(apperr.CodeBadRequest, "patient has no appointment history", nil)
	}

	did := utils.StringToUUIDv7(appointment.DoctorID)
	order := &models.Order{
		ID:        utils.GenerateUUIDv7(),
		PatientID: patientID,
		DoctorID:  &did,
		Note:      body.Note,
		Status:    models.OrderStatusPending,
	}

	if err := s.orderRepository.Create(ctx, order); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to create order", err)
	}

	return &dto.CreateOrderResponseDto{
		OrderID: order.ID.String(),
	}, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, body dto.UpdateOrderRequestDto) (*dto.UpdateOrderResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "doctor" {
		return nil, apperr.New(apperr.CodeForbidden, "only doctors can update orders", nil)
	}
	orderID := body.OrderID
	parsedOrderID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid order ID", err)
	}

	order, err := s.orderRepository.FindByID(ctx, parsedOrderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "order not found", err)
	}

	doctorID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}

	if order.DoctorID == nil || *order.DoctorID != doctorID {
		return nil, apperr.New(apperr.CodeForbidden, "doctor can only edit their own orders", nil)
	}

	// main logic begins here
	// delete existing order items
	if err := s.orderItemRepository.DeleteByOrderID(ctx, order.ID); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to delete existing order items", err)
	}
	// add new order items
	for _, item := range body.OrderItems {
		medicine, err := s.medicineRepository.FindByID(ctx, item.MedicineID)
		if err != nil {
			return nil, apperr.New(apperr.CodeBadRequest, "medicine not found", err)
		}
		orderItem := &models.OrderItem{
			ID:         utils.GenerateUUIDv7(),
			OrderID:    order.ID,
			MedicineID: medicine.ID,
			Quantity:   item.Quantity,
		}
		if err := s.orderItemRepository.Create(ctx, orderItem); err != nil {
			return nil, apperr.New(apperr.CodeInternal, "failed to create order item", err)
		}
	}

	return &dto.UpdateOrderResponseDto{
		OrderID: order.ID.String(),
	}, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID string) (*dto.GetOrderByIDResponseDto, error) {
	parsedOrderID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid order ID", err)
	}

	order, err := s.orderRepository.FindByID(ctx, parsedOrderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "order not found", err)
	}

	// Convert order items to response format
	orderItems := make([]dto.OrderItem, len(order.OrderItems))
	for i, item := range order.OrderItems {
		medicineName := ""
		if item.Medicine != nil {
			medicineName = item.Medicine.Name
		}
		orderItems[i] = dto.OrderItem{
			MedicineID:   item.MedicineID.String(),
			MedicineName: medicineName,
			Quantity:     item.Quantity,
		}
	}

	// Format timestamps
	var submittedAt, reviewedAt *string
	if order.SubmittedAt != nil {
		submittedAtStr := order.SubmittedAt.Format("2006-01-02T15:04:05Z07:00")
		submittedAt = &submittedAtStr
	}
	if order.ReviewedAt != nil {
		reviewedAtStr := order.ReviewedAt.Format("2006-01-02T15:04:05Z07:00")
		reviewedAt = &reviewedAtStr
	}

	// Fetch delivery information if exists
	var deliveryStatus, deliveryAt *string
	delivery, err := s.deliveryRepository.FindByOrderID(ctx, parsedOrderID)
	if err == nil && delivery != nil {
		status := string(delivery.Status)
		deliveryStatus = &status
		if delivery.DeliveredAt != nil {
			deliveredAtStr := delivery.DeliveredAt.Format("2006-01-02T15:04:05Z07:00")
			deliveryAt = &deliveredAtStr
		}
	}

	return &dto.GetOrderByIDResponseDto{
		OrderID:        order.ID.String(),
		PatientID:      order.PatientID.String(),
		DoctorID:       order.DoctorID.String(),
		TotalAmount:    order.TotalAmount,
		Note:           order.Note,
		SubmittedAt:    submittedAt,
		ReviewedAt:     reviewedAt,
		Status:         string(order.Status),
		DeliveryStatus: deliveryStatus,
		DeliveryAt:     deliveryAt,
		OrderItems:     orderItems,
	}, nil
}

func (s *OrderService) GetAllOrdersHistoryByPatientID(ctx context.Context) (*dto.GetAllOrdersHistoryListDto, error) {
	userID := contextUtils.GetUserId(ctx)

	patientID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}

	orders, err := s.orderRepository.FindByPatientID(ctx, patientID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve orders", err)
	}

	orderHistoryList := make([]dto.GetAllOrdersHistoryResponseDto, len(orders))

	for idx, order := range orders {
		// Convert order items to response format
		orderItems := make([]dto.OrderItem, len(order.OrderItems))
		for i, item := range order.OrderItems {
			medicineName := ""
			if item.Medicine != nil {
				medicineName = item.Medicine.Name
			}
			orderItems[i] = dto.OrderItem{
				MedicineID:   item.MedicineID.String(),
				MedicineName: medicineName,
				Quantity:     item.Quantity,
			}
		}

		// Format timestamps
		var submittedAt, reviewedAt *string
		if order.SubmittedAt != nil {
			submittedAtStr := order.SubmittedAt.Format("2006-01-02T15:04:05Z07:00")
			submittedAt = &submittedAtStr
		}
		if order.ReviewedAt != nil {
			reviewedAtStr := order.ReviewedAt.Format("2006-01-02T15:04:05Z07:00")
			reviewedAt = &reviewedAtStr
		}

		// Fetch delivery information if exists
		var deliveryStatus, deliveryAt *string
		delivery, err := s.deliveryRepository.FindByOrderID(ctx, order.ID)
		if err == nil && delivery != nil {
			status := string(delivery.Status)
			deliveryStatus = &status
			if delivery.DeliveredAt != nil {
				deliveredAtStr := delivery.DeliveredAt.Format("2006-01-02T15:04:05Z07:00")
				deliveryAt = &deliveredAtStr
			}
		}

		var doctorID *string
		if order.DoctorID != nil {
			doctorIDStr := order.DoctorID.String()
			doctorID = &doctorIDStr
		}

		orderHistoryList[idx] = dto.GetAllOrdersHistoryResponseDto{
			OrderID:        order.ID.String(),
			PatientID:      order.PatientID.String(),
			DoctorID:       doctorID,
			TotalAmount:    order.TotalAmount,
			Note:           order.Note,
			SubmittedAt:    submittedAt,
			ReviewedAt:     reviewedAt,
			Status:         string(order.Status),
			DeliveryStatus: deliveryStatus,
			DeliveryAt:     deliveryAt,
			CreatedAt:      order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			OrderItems:     orderItems,
		}
	}

	return &dto.GetAllOrdersHistoryListDto{
		Orders: orderHistoryList,
		Total:  len(orderHistoryList),
	}, nil
}

func (s *OrderService) GetLatestOrderByPatientID(ctx context.Context) (*dto.GetOrderByIDResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)

	patientID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}

	order, err := s.orderRepository.FindLatestOrderByPatientID(ctx, patientID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve order", err)
	}
	if order == nil {
		return nil, apperr.New(apperr.CodeNotFound, "no orders found", nil)
	}

	// Convert order items to response format
	orderItems := make([]dto.OrderItem, len(order.OrderItems))
	for i, item := range order.OrderItems {
		medicineName := ""
		if item.Medicine != nil {
			medicineName = item.Medicine.Name
		}
		orderItems[i] = dto.OrderItem{
			MedicineID:   item.MedicineID.String(),
			MedicineName: medicineName,
			Quantity:     item.Quantity,
		}
	}

	// Format timestamps
	var submittedAt, reviewedAt *string
	if order.SubmittedAt != nil {
		submittedAtStr := order.SubmittedAt.Format("2006-01-02T15:04:05Z07:00")
		submittedAt = &submittedAtStr
	}
	if order.ReviewedAt != nil {
		reviewedAtStr := order.ReviewedAt.Format("2006-01-02T15:04:05Z07:00")
		reviewedAt = &reviewedAtStr
	}

	// Fetch delivery information if exists
	var deliveryStatus, deliveryAt *string
	delivery, err := s.deliveryRepository.FindByOrderID(ctx, order.ID)
	if err == nil && delivery != nil {
		status := string(delivery.Status)
		deliveryStatus = &status
		if delivery.DeliveredAt != nil {
			deliveredAtStr := delivery.DeliveredAt.Format("2006-01-02T15:04:05Z07:00")
			deliveryAt = &deliveredAtStr
		}
	}

	return &dto.GetOrderByIDResponseDto{
		OrderID:        order.ID.String(),
		PatientID:      order.PatientID.String(),
		DoctorID:       order.DoctorID.String(),
		TotalAmount:    order.TotalAmount,
		Note:           order.Note,
		SubmittedAt:    submittedAt,
		ReviewedAt:     reviewedAt,
		Status:         string(order.Status),
		DeliveryStatus: deliveryStatus,
		DeliveryAt:     deliveryAt,
		OrderItems:     orderItems,
	}, nil
}

func (s *OrderService) GetLatestOrderByPatientIDForDoctor(ctx context.Context, patientID string) (*dto.GetOrderByIDResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "doctor" {
		return nil, apperr.New(apperr.CodeForbidden, "only doctors can access this endpoint", nil)
	}

	parsedPatientID, err := uuid.Parse(patientID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid patient ID", err)
	}

	order, err := s.orderRepository.FindLatestOrderByPatientID(ctx, parsedPatientID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve order", err)
	}
	if order == nil {
		return nil, apperr.New(apperr.CodeNotFound, "no orders found for this patient", nil)
	}

	// Verify the doctor is the one assigned to the order
	doctorID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}

	if order.DoctorID == nil || *order.DoctorID != doctorID {
		return nil, apperr.New(apperr.CodeForbidden, "doctor can only access their own patient's orders", nil)
	}

	// Convert order items to response format
	orderItems := make([]dto.OrderItem, len(order.OrderItems))
	for i, item := range order.OrderItems {
		medicineName := ""
		if item.Medicine != nil {
			medicineName = item.Medicine.Name
		}
		orderItems[i] = dto.OrderItem{
			MedicineID:   item.MedicineID.String(),
			MedicineName: medicineName,
			Quantity:     item.Quantity,
		}
	}

	// Format timestamps
	var submittedAt, reviewedAt *string
	if order.SubmittedAt != nil {
		submittedAtStr := order.SubmittedAt.Format("2006-01-02T15:04:05Z07:00")
		submittedAt = &submittedAtStr
	}
	if order.ReviewedAt != nil {
		reviewedAtStr := order.ReviewedAt.Format("2006-01-02T15:04:05Z07:00")
		reviewedAt = &reviewedAtStr
	}

	// Fetch delivery information if exists
	var deliveryStatus, deliveryAt *string
	delivery, err := s.deliveryRepository.FindByOrderID(ctx, order.ID)
	if err == nil && delivery != nil {
		status := string(delivery.Status)
		deliveryStatus = &status
		if delivery.DeliveredAt != nil {
			deliveredAtStr := delivery.DeliveredAt.Format("2006-01-02T15:04:05Z07:00")
			deliveryAt = &deliveredAtStr
		}
	}

	return &dto.GetOrderByIDResponseDto{
		OrderID:        order.ID.String(),
		PatientID:      order.PatientID.String(),
		DoctorID:       order.DoctorID.String(),
		TotalAmount:    order.TotalAmount,
		Note:           order.Note,
		SubmittedAt:    submittedAt,
		ReviewedAt:     reviewedAt,
		Status:         string(order.Status),
		DeliveryStatus: deliveryStatus,
		DeliveryAt:     deliveryAt,
		OrderItems:     orderItems,
	}, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, body dto.CancelOrderRequestDto) (*dto.CancelOrderResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "doctor" {
		return nil, apperr.New(apperr.CodeForbidden, "only doctors can cancel orders", nil)
	}

	parsedOrderID, err := uuid.Parse(body.OrderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid order ID", err)
	}

	order, err := s.orderRepository.FindByID(ctx, parsedOrderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "order not found", err)
	}

	doctorID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}

	if order.DoctorID == nil || *order.DoctorID != doctorID {
		return nil, apperr.New(apperr.CodeForbidden, "doctor can only cancel their own orders", nil)
	}

	order.Status = models.OrderStatusCancelled
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to cancel order", err)
	}

	return &dto.CancelOrderResponseDto{
		OrderID: order.ID.String(),
		Status:  string(order.Status),
	}, nil
}

func (s *OrderService) ApproveOrder(ctx context.Context, body dto.ApproveOrderRequestDto) (*dto.ApproveOrderResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "doctor" {
		return nil, apperr.New(apperr.CodeForbidden, "only doctors can approve orders", nil)
	}

	parsedOrderID, err := uuid.Parse(body.OrderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid order ID", err)
	}

	order, err := s.orderRepository.FindByID(ctx, parsedOrderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "order not found", err)
	}

	doctorID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}

	if order.DoctorID == nil || *order.DoctorID != doctorID {
		return nil, apperr.New(apperr.CodeForbidden, "doctor can only approve their own orders", nil)
	}

	order.Status = models.OrderStatusApproved
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to approve order", err)
	}

	return &dto.ApproveOrderResponseDto{
		OrderID: order.ID.String(),
		Status:  string(order.Status),
	}, nil
}
