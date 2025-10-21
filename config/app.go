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
    })

    app.Use(logger.New())
    app.Use(cors.New())

    return app
}