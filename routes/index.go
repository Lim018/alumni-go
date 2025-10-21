package routes

import (
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *mongo.Database) {
    AlumniRoutes(app, db)
    PekerjaanRoutes(app, db)
    AuthRoutes(app, db) 
    UserRoutes(app, db)
}