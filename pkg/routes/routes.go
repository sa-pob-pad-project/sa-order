package routes

import (
	// "user-service/pkg/context"
	// "user-service/pkg/dto"
	_ "order-service/docs"
	"order-service/pkg/handlers"
	"order-service/pkg/jwt"
	"order-service/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, orderHandler *handlers.OrderHandler, deliveryInfoHandler *handlers.DeliveryInfoHandler, jwtSvc *jwt.JwtService) {

	api := app.Group("/api")
	api.Get("/swagger/*", swagger.HandlerDefault)

	// Order Routes
	order := api.Group("/order")
	orderV1 := order.Group("/v1")
	orderV1.Use(middleware.JwtMiddleware(jwtSvc))
	orderV1.Post("/orders", orderHandler.CreateOrder)
	orderV1.Delete("/orders", orderHandler.CancelOrder)
	orderV1.Get("/orders/:id", orderHandler.GetOrder)
	orderV1.Get("/orders", orderHandler.GetAllOrdersHistory)
	// orderV1.Put("/orders", )
	// orderV1.Post("/orders/confirm", )
	// orderV1.Get("/orders/latest", )
	// orderV1.Get("/orders/latest/:patient_id", )

	// Delivery Routes
	delivery := api.Group("/delivery")
	deliveryV1 := delivery.Group("/v1")
	deliveryV1.Use(middleware.JwtMiddleware(jwtSvc))

	// Delivery Information Routes
	deliveryInfo := api.Group("/delivery-info")
	deliveryInfoV1 := deliveryInfo.Group("/v1")
	deliveryInfoV1.Use(middleware.JwtMiddleware(jwtSvc))
}
