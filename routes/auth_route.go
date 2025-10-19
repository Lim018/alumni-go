package routes

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/service"
	"go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, db *sql.DB) {
	app.Post("/login", func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Request body tidak valid",
			})
		}

		resp, err := service.LoginService(db, req)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		// Bungkus response
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Login berhasil",
			"data": fiber.Map{
				"user": fiber.Map{
					"id":       resp.User.ID,
					"username": resp.User.Username,
					"email":    resp.User.Email,
					"role":     resp.User.Role,
				},
				"token": resp.Token,
			},
		})
	})
}

func UserRoutes(app *fiber.App, db *sql.DB) {
	users := app.Group("/users", middleware.AuthRequired())
	users.Get("/", service.GetUsersService(db))
}
