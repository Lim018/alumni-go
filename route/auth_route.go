package route

import (
	"alumni-go/app/model"
	"alumni-go/app/service"
	"alumni-go/helper"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler sekarang menerima db dan jwtSecret
func NewAuthHandler(db *sql.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		// Teruskan kedua dependensi saat membuat service
		authService: service.NewAuthService(db, jwtSecret),
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return helper.SuccessResponse(c, "Login successful", response)
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return helper.SuccessResponse(c, "Profile retrieved successfully", user)
}
