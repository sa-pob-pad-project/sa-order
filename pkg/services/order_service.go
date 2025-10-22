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
	"time"

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

// calculateOrderTotal calculates the total amount for an order based on its items
func (s *OrderService) calculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	orderItems, err := s.orderItemRepository.FindByOrderID(ctx, orderID)
	if err != nil {
		return 0, err
	}

	var totalAmount float64
	for _, item := range orderItems {
		if item.Medicine != nil {
			totalAmount += item.Medicine.Price * item.Quantity
		}
	}

	return totalAmount, nil
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
	submittedAt := time.Now()
	order := &models.Order{
		ID:          utils.GenerateUUIDv7(),
		PatientID:   patientID,
		DoctorID:    &did,
		Note:        body.Note,
		Status:      models.OrderStatusPending,
		SubmittedAt: &submittedAt,
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
	var totalAmount float64
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
		// Calculate total amount
		totalAmount += medicine.Price * item.Quantity
	}

	// Update order with calculated total amount
	order.TotalAmount = totalAmount
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to update order total amount", err)
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
		return &dto.GetOrderByIDResponseDto{}, nil
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
	if order == nil { // return empty response if no orders found
		return nil, nil
	}

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

	// Calculate total amount before approving
	totalAmount, err := s.calculateOrderTotal(ctx, order.ID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to calculate order total", err)
	}

	order.Status = models.OrderStatusApproved
	order.TotalAmount = totalAmount
	reviewedAt := time.Now()
	order.ReviewedAt = &reviewedAt
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to approve order", err)
	}

	return &dto.ApproveOrderResponseDto{
		OrderID: order.ID.String(),
		Status:  string(order.Status),
	}, nil
}

func (s *OrderService) RejectOrder(ctx context.Context, body dto.RejectOrderRequestDto) (*dto.RejectOrderResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "doctor" {
		return nil, apperr.New(apperr.CodeForbidden, "only doctors can reject orders", nil)
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
		return nil, apperr.New(apperr.CodeForbidden, "doctor can only reject their own orders", nil)
	}

	// Calculate total amount before rejecting
	totalAmount, err := s.calculateOrderTotal(ctx, order.ID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to calculate order total", err)
	}

	order.Status = models.OrderStatusRejected
	order.TotalAmount = totalAmount
	reviewedAt := time.Now()
	order.ReviewedAt = &reviewedAt
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to reject order", err)
	}

	return &dto.RejectOrderResponseDto{
		OrderID: order.ID.String(),
		Status:  string(order.Status),
	}, nil
}

func (s *OrderService) PayOrder(ctx context.Context, body dto.PayOrderRequestDto) (*dto.PayOrderResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	// อนุญาตเฉพาะคนไข้/ผู้สั่งซื้อ (ปรับตามระบบคุณ: "patient", "user" หรือ role อื่น)
	if role != "patient" {
		return nil, apperr.New(apperr.CodeForbidden, "only patients can pay orders", nil)
	}

	// validate body
	if body.OrderID == "" {
		return nil, apperr.New(apperr.CodeBadRequest, "order ID is required", nil)
	}
	parsedOrderID, err := uuid.Parse(body.OrderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid order ID", err)
	}

	// โหลดออเดอร์
	order, err := s.orderRepository.FindByID(ctx, parsedOrderID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "order not found", err)
	}

	patientID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}
	if order.PatientID != patientID {
		return nil, apperr.New(apperr.CodeForbidden, "patient can only pay their own orders", nil)
	}

	// ตรวจสถานะที่อนุญาตให้จ่ายเงิน
	switch order.Status {
	case models.OrderStatusPaid, models.OrderStatusProcessing, models.OrderStatusShipped, models.OrderStatusDelivered:
		return nil, apperr.New(apperr.CodeConflict, "order already paid or processed", nil)
	case models.OrderStatusCancelled, models.OrderStatusRejected:
		return nil, apperr.New(apperr.CodeForbidden, "order cannot be paid in current state", nil)
	case models.OrderStatusPending, models.OrderStatusApproved:
		// allowed
	default:
		return nil, apperr.New(apperr.CodeBadRequest, "unknown order state", nil)
	}

	// คำนวณยอดรวมล่าสุด (กันกรณีมีส่วนลด/ราคาเปลี่ยน)
	totalAmount, err := s.calculateOrderTotal(ctx, order.ID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to calculate order total", err)
	}

	order.Status = models.OrderStatusPaid
	order.TotalAmount = totalAmount

	if err := s.orderRepository.Update(ctx, order); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to mark order paid", err)
	}

	return &dto.PayOrderResponseDto{
		OrderID: order.ID.String(),
		Status:  string(order.Status),
	}, nil
}

func (s *OrderService) GetAllOrdersByDoctorID(ctx context.Context) (*dto.GetAllOrdersForDoctorListDto, error) {
	userID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "doctor" {
		return nil, apperr.New(apperr.CodeForbidden, "only doctors can access this endpoint", nil)
	}

	doctorID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}

	orders, err := s.orderRepository.FindByDoctorID(ctx, doctorID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve orders", err)
	}

	// Collect all unique patient IDs
	patientIDMap := make(map[string]bool)
	patientIDs := []string{}
	for _, order := range orders {
		patientID := order.PatientID.String()
		if !patientIDMap[patientID] {
			patientIDMap[patientID] = true
			patientIDs = append(patientIDs, patientID)
		}
	}

	// Fetch patient profiles from user client
	patientProfiles := make(map[string]*dto.PatientInfo)
	if len(patientIDs) > 0 {
		profiles, err := s.userClient.GetPatientByIds(ctx, patientIDs)
		if err == nil && profiles != nil {
			for _, profile := range *profiles {
				patientProfiles[profile.ID] = &dto.PatientInfo{
					PatientID:   profile.ID,
					FirstName:   profile.FirstName,
					LastName:    profile.LastName,
					Gender:      profile.Gender,
					PhoneNumber: profile.PhoneNumber,
				}
			}
		}
	}

	orderHistoryList := make([]dto.GetAllOrdersForDoctorResponseDto, len(orders))

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

		var doctorIDStr *string
		if order.DoctorID != nil {
			doctorIDStrVal := order.DoctorID.String()
			doctorIDStr = &doctorIDStrVal
		}

		// Get patient info from the map
		patientInfo := patientProfiles[order.PatientID.String()]

		orderHistoryList[idx] = dto.GetAllOrdersForDoctorResponseDto{
			OrderID:        order.ID.String(),
			PatientID:      order.PatientID.String(),
			PatientInfo:    patientInfo,
			DoctorID:       doctorIDStr,
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

	return &dto.GetAllOrdersForDoctorListDto{
		Orders: orderHistoryList,
		Total:  len(orderHistoryList),
	}, nil
}

func (s *OrderService) GetAllOrdersHistoryForDoctor(ctx context.Context, statusFilter string) (*dto.GetAllOrdersForDoctorListDto, error) {
	userID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "doctor" {
		return nil, apperr.New(apperr.CodeForbidden, "only doctors can access this endpoint", nil)
	}

	doctorID, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
	}

	var orders []models.Order

	// Filter by status if provided (approved or rejected)
	if statusFilter != "" {
		status := models.OrderStatus(statusFilter)
		orders, err = s.orderRepository.FindByDoctorIDAndStatus(ctx, doctorID, status)
		if err != nil {
			return nil, apperr.New(apperr.CodeInternal, "failed to retrieve orders", err)
		}
	} else {
		// If no filter, get all non-pending orders (approved or rejected)
		allOrders, err := s.orderRepository.FindByDoctorIDAndStatuses(ctx, doctorID, []models.OrderStatus{models.OrderStatusApproved, models.OrderStatusRejected})
		if err != nil {
			return nil, apperr.New(apperr.CodeInternal, "failed to retrieve orders", err)
		}

		orders = allOrders
	}

	// Collect all unique patient IDs
	patientIDMap := make(map[string]bool)
	patientIDs := []string{}
	for _, order := range orders {
		patientID := order.PatientID.String()
		if !patientIDMap[patientID] {
			patientIDMap[patientID] = true
			patientIDs = append(patientIDs, patientID)
		}
	}

	// Fetch patient profiles from user client
	patientProfiles := make(map[string]*dto.PatientInfo)
	if len(patientIDs) > 0 {
		profiles, err := s.userClient.GetPatientByIds(ctx, patientIDs)
		if err == nil && profiles != nil {
			for _, profile := range *profiles {
				patientProfiles[profile.ID] = &dto.PatientInfo{
					PatientID:   profile.ID,
					FirstName:   profile.FirstName,
					LastName:    profile.LastName,
					Gender:      profile.Gender,
					PhoneNumber: profile.PhoneNumber,
				}
			}
		}
	}

	orderHistoryList := make([]dto.GetAllOrdersForDoctorResponseDto, len(orders))

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

		var doctorIDStr *string
		if order.DoctorID != nil {
			doctorIDStrVal := order.DoctorID.String()
			doctorIDStr = &doctorIDStrVal
		}

		// Get patient info from the map
		patientInfo := patientProfiles[order.PatientID.String()]

		orderHistoryList[idx] = dto.GetAllOrdersForDoctorResponseDto{
			OrderID:        order.ID.String(),
			PatientID:      order.PatientID.String(),
			PatientInfo:    patientInfo,
			DoctorID:       doctorIDStr,
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

	return &dto.GetAllOrdersForDoctorListDto{
		Orders: orderHistoryList,
		Total:  len(orderHistoryList),
	}, nil
}
