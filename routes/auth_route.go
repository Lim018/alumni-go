package routes

import (
    "go-fiber/app/model"
    "go-fiber/app/service"
    
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

func AuthRoutes(app *fiber.App, db *mongo.Database) {
    auth := app.Group("/auth")

    // Login godoc
    // @Summary User login
    // @Description Authenticate user and return JWT token
    // @Tags Authentication
    // @Accept json
    // @Produce json
    // @Param body body model.LoginRequest true "Login credentials"
    // @Success 200 {object} map[string]interface{} "Login successful"
    // @Failure 400 {object} map[string]interface{} "Invalid request"
    // @Failure 401 {object} map[string]interface{} "Unauthorized"
    // @Router /auth/login [post]
    auth.Post("/login", func(c *fiber.Ctx) error {
        var req model.LoginRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{
                "error":   "Invalid request",
                "success": false,
            })
        }

        response, err := service.LoginService(db, req)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error":   err.Error(),
                "success": false,
            })
        }

        return c.JSON(fiber.Map{
            "message": "Login successful",
            "success": true,
            "data":    response,
        })
    })
}