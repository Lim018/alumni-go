package middleware

import (
    "strings"
    
    "go-fiber/utils"
    
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error":   "Missing authorization header",
                "success": false,
            })
        }

        // Extract token from "Bearer <token>"
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            return c.Status(401).JSON(fiber.Map{
                "error":   "Invalid authorization format",
                "success": false,
            })
        }

        claims, err := utils.ParseToken(tokenString)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error":   "Invalid or expired token",
                "success": false,
            })
        }

        // Convert string UserID to ObjectID
        userID, err := primitive.ObjectIDFromHex(claims.UserID)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error":   "Invalid user ID in token",
                "success": false,
            })
        }

        // Store user info in context
        c.Locals("user_id", userID)
        c.Locals("username", claims.Username)
        c.Locals("role", claims.Role)

        return c.Next()
    }
}

func AdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        role := c.Locals("role")
        if role == nil || role.(string) != "admin" {
            return c.Status(403).JSON(fiber.Map{
                "error":   "Access denied. Admin only.",
                "success": false,
            })
        }
        return c.Next()
    }
}