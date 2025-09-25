package route

import (
	"alumni-go/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Initialize handlers
	authHandler := NewAuthHandler()
	alumniHandler := NewAlumniHandler()
	pekerjaanHandler := NewPekerjaanHandler()

	// API group
	api := app.Group("/alumni-management-system")

	// Auth routes (no middleware)
	api.Post("/login", authHandler.Login)

	// Protected routes
	protected := api.Use(middleware.JWTMiddleware())
	
	// Profile route (for both admin and user)
	protected.Get("/profile", middleware.AdminOrUser(), authHandler.Profile)

	// Alumni routes
	alumni := protected.Group("/alumni")
	alumni.Get("/", middleware.AdminOrUser(), alumniHandler.GetAll)          // Admin + User
	alumni.Get("/:id", middleware.AdminOrUser(), alumniHandler.GetByID)      // Admin + User
	alumni.Post("/", middleware.AdminOnly(), alumniHandler.Create)           // Admin only
	alumni.Put("/:id", middleware.AdminOnly(), alumniHandler.Update)         // Admin only
	alumni.Delete("/:id", middleware.AdminOnly(), alumniHandler.Delete)      // Admin only

	alumni.Get("/laporan/baru-bekerja", middleware.AdminOrUser(), alumniHandler.GetAlumniBaruBekerja)

	// Pekerjaan routes
	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", middleware.AdminOrUser(), pekerjaanHandler.GetAll)                    // Admin + User
	pekerjaan.Get("/:id", middleware.AdminOrUser(), pekerjaanHandler.GetByID)                // Admin + User
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), pekerjaanHandler.GetByAlumniID) // Admin only
	pekerjaan.Post("/", middleware.AdminOnly(), pekerjaanHandler.Create)                     // Admin only
	pekerjaan.Put("/:id", middleware.AdminOnly(), pekerjaanHandler.Update)                   // Admin only
	pekerjaan.Delete("/:id", middleware.AdminOnly(), pekerjaanHandler.Delete)                // Admin only

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
			"message": "Alumni Management System is running",
		})
	})
}