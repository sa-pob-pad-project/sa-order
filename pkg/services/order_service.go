package service

import (
	"order-service/pkg/clients"
	"order-service/pkg/repository"

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
