package middleware

import (
	"order-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func JwtMiddleware(jwtService *jwt.JwtService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := c.Cookies("access_token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		claims, err := jwtService.Parse(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("accessToken", token)

		return c.Next()
	}
}
