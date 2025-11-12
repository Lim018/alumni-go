package routes

import (
    "go-fiber/app/service"
    "go-fiber/middleware"
    
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

func FileRoutes(app *fiber.App, db *mongo.Database) {
    files := app.Group("/files", middleware.AuthRequired())

    files.Post("/upload/foto", func(c *fiber.Ctx) error {
        return service.UploadFotoService(c, db)
    })

    files.Post("/upload/sertifikat", func(c *fiber.Ctx) error {
        return service.UploadSertifikatService(c, db)
    })

    // Get files by alumni
    files.Get("/alumni/:alumni_id", func(c *fiber.Ctx) error {
        return service.GetFilesByAlumniService(c, db)
    })

    // Get all files (with pagination) - Admin only
    files.Get("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
        return service.GetAllFilesService(c, db)
    })

    // Download file
    files.Get("/:id/download", func(c *fiber.Ctx) error {
        return service.DownloadFileService(c, db)
    })

    // Delete file
    files.Delete("/:id", func(c *fiber.Ctx) error {
        return service.DeleteFileService(c, db)
    })
}