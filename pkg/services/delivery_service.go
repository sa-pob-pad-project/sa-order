package service

import (
	"order-service/pkg/clients"
	"order-service/pkg/repository"

	"gorm.io/gorm"
)

type DeliveryService struct {
	db                     *gorm.DB
	deliveryRepository     *repository.DeliveryRepository
	deliveryInfoRepository *repository.DeliveryInformationRepository
	orderRepository        *repository.OrderRepository
	userClient             *clients.UserClient
}

func NewDeliveryService(
	db *gorm.DB,
	deliveryRepo *repository.DeliveryRepository,
	deliveryInfoRepo *repository.DeliveryInformationRepository,
	orderRepo *repository.OrderRepository,
	userClient *clients.UserClient,
) *DeliveryService {
	return &DeliveryService{
		db:                     db,
		deliveryRepository:     deliveryRepo,
		deliveryInfoRepository: deliveryInfoRepo,
		orderRepository:        orderRepo,
		userClient:             userClient,
	}
}
