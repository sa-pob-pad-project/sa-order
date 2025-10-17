package contextUtils

import (
	"context"
	// "time"
	"order-service/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type contextKey string

const (
	ContextKeyUserID      contextKey = "userID"
	ContextKeyRole        contextKey = "role"
	ContextKeyAccessToken contextKey = "accessToken"
)

func WithBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		if err := c.BodyParser(&body); err != nil {
			return response.BadRequest(c, "Cannot parse JSON")
		}
		validate := validator.New()
		if err := validate.Struct(body); err != nil {
			return response.BadRequest(c, err.Error())
		}
		c.Locals("body", body)
		return c.Next()
	}
}

func GetUserId(c context.Context) string {
	return c.Value(ContextKeyUserID).(string)
}

func GetRole(c context.Context) string {
	return c.Value(ContextKeyRole).(string)
}

func GetAccessToken(c context.Context) string {
	return c.Value(ContextKeyAccessToken).(string)
}

func GetContext(c *fiber.Ctx) context.Context {
	ctx := c.UserContext()
	userID := c.Locals("userID")
	if s, ok := userID.(string); ok {
		ctx = context.WithValue(ctx, ContextKeyUserID, s)
	}
	role := c.Locals("role")
	if s, ok := role.(string); ok {
		ctx = context.WithValue(ctx, ContextKeyRole, s)
	}
	token := c.Locals("accessToken")
	if s, ok := token.(string); ok {
		ctx = context.WithValue(ctx, ContextKeyAccessToken, s)
	}

	return ctx
}
