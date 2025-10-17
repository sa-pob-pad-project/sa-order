package handlers

import (
	service "order-service/pkg/services"
)

type DeliveryInfoHandler struct {
	deliveryService *service.DeliveryService
}

func NewDeliveryInfoHandler(deliveryService *service.DeliveryService) *DeliveryInfoHandler {
	return &DeliveryInfoHandler{
		deliveryService: deliveryService,
	}
}
