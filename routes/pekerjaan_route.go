package routes

import (
	"database/sql"
	"go-fiber/app/service"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoutes(app *fiber.App, db *sql.DB) {
	pekerjaan := app.Group("/pekerjaan", middleware.AuthRequired())

    pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
        return service.GetPekerjaanByAlumniIDService(c, db)
    })

	pekerjaan.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanService(c, db)
	})

	pekerjaan.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanService(c, db)
	})

	pekerjaan.Delete("/:id", func(c *fiber.Ctx) error {
    return service.SoftDeletePekerjaanService(c, db)
	})

	pekerjaan.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanServiceDatatable(c, db)
	})

	trash := pekerjaan.Group("/trash")

	trash.Get("/", func(c *fiber.Ctx) error {
		return service.GetTrashPekerjaanService(c, db)
	})

	trash.Put("/restore/:id", func(c *fiber.Ctx) error {
		return service.RestorePekerjaanService(c, db)
	})

	trash.Delete("/:id", func(c *fiber.Ctx) error {
		return service.HardDeletePekerjaanService(c, db)
	})
}