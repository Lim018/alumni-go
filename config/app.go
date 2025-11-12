package config

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "go.mongodb.org/mongo-driver/mongo"
)

func NewApp(db *mongo.Database) *fiber.App {
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "error":   err.Error(),
                "success": false,
            })
        },
        BodyLimit: 2 * 1024 * 1024, // 2MB max body size (for sertifikat)
    })

    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

    // Serve static files
    app.Static("/uploads/foto", "./uploads/foto")
    app.Static("/uploads/sertifikat", "./uploads/sertifikat")

    return app
}