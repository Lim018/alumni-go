package routes

import (
    "go-fiber/app/service"
    "go-fiber/middleware"
    
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

func AlumniRoutes(app *fiber.App, db *mongo.Database) {
    alumni := app.Group("/alumni", middleware.AuthRequired())

    alumni.Get("/", func(c *fiber.Ctx) error {
        return service.GetAllAlumniServiceDatatable(c, db)
    })

    alumni.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
        return service.CreateAlumniService(c, db)
    })

    alumni.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
        return service.UpdateAlumniService(c, db)
    })

    alumni.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
        return service.DeleteAlumniService(c, db)
    })

    alumni.Get("/stats/jurusan", func(c *fiber.Ctx) error {
        return service.GetAlumniStatsService(c, db)
    })
}