package routes

import (
	// "user-service/pkg/context"
	// "user-service/pkg/dto"
	_ "order-service/docs"
	"order-service/pkg/handlers"
	"order-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, orderHandler *handlers.OrderHandler, deliveryInfoHandler *handlers.DeliveryInfoHandler, jwtSvc *jwt.JwtService) {

	api := app.Group("/api")
	api.Get("/swagger/*", swagger.HandlerDefault)
	// Order Routes

	// Delivery Routes

	// Delivery Information Routes
}
