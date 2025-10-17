package handlers

import (
	service "order-service/pkg/services"
)

type OrderHandler struct {
	orderService    *service.OrderService
	deliveryService *service.DeliveryService
}

func NewOrderHandler(orderService *service.OrderService, deliveryService *service.DeliveryService) *OrderHandler {
	return &OrderHandler{
		orderService:    orderService,
		deliveryService: deliveryService,
	}
}
