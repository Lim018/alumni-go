package helper

import (
	"alumni-go/app/model"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(model.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(model.APIResponse{
		Success: false,
		Message: message,
	})
}

func PaginatedSuccessResponse(c *fiber.Ctx, message string, data interface{}, meta model.MetaData) error {
	return c.JSON(model.PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func CalculateTotalPages(total int64, perPage int) int {
	if perPage == 0 {
		return 0
	}
	return int((total + int64(perPage) - 1) / int64(perPage))
}