package routes

import (
    "go-fiber/app/service"
    "go-fiber/middleware"
    
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(app *fiber.App, db *mongo.Database) {
    users := app.Group("/users", middleware.AuthRequired(), middleware.AdminOnly())

    users.Get("/", service.GetUsersService(db))
}