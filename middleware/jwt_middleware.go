package middleware

import (
	"alumni-management-system/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		auth := c.Get("Authorization")
		if auth == "" {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization token")
		}

		// Check if it's a Bearer token
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization format")
		}

		token := parts[1]

		// Validate token
		claims, err := helper.ValidateToken(token)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Set user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}