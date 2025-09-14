package route

import (
	"alumni-management-system/app/model"
	"alumni-management-system/app/service"
	"alumni-management-system/helper"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Basic validation
	if req.Username == "" || req.Password == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Username and password are required")
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