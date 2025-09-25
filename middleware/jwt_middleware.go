package middleware

import (
	"alumni-go/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Struct untuk menampung dependensi middleware, seperti secret key
type Middleware struct {
	jwtSecret string
}

// Constructor untuk membuat instance Middleware
func NewMiddleware(jwtSecret string) *Middleware {
	return &Middleware{jwtSecret: jwtSecret}
}

// JWTMiddleware sekarang menjadi method dari struct Middleware
func (m *Middleware) JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization token")
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization format")
		}

		token := parts[1]

		// Validasi token menggunakan secret key dari struct
		claims, err := helper.ValidateToken(token, m.jwtSecret)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
