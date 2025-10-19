package routes

import (
	"database/sql"
	"go-fiber/app/service"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func AlumniRoutes(app *fiber.App, db *sql.DB) {
	alumni := app.Group("/alumni", middleware.AuthRequired()) 


	alumni.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	})

	alumni.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})

	alumni.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})

	alumni.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllAlumniServiceDatatable(c, db)
	})

	alumni.Get("/stats/jurusan", func(c *fiber.Ctx) error {
        return service.GetAlumniStatsService(c, db)
    })
}
