package response

import "github.com/gofiber/fiber/v2"

type BaseResponse struct {
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func OK[T any](c *fiber.Ctx, data T) error {
	return c.Status(fiber.StatusOK).JSON(data)
}

func Created[T any](c *fiber.Ctx, data T) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": message,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": message,
	})
}

func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": message,
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": message,
	})
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": message,
	})
}

func Failed(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}
