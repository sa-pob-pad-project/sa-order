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
