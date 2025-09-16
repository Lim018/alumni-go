package middleware

import (
	"alumni-go/helper"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "User role not found")
		}

		// Check if user role is in allowed roles
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return helper.ErrorResponse(c, fiber.StatusForbidden, "Insufficient permissions")
	}
}

func AdminOnly() fiber.Handler {
	return RequireRole("admin")
}

func AdminOrUser() fiber.Handler {
	return RequireRole("admin", "user")
}