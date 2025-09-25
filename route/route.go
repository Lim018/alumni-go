package route

import (
	"alumni-go/middleware"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes menerima db dan jwtSecret untuk di-inject ke lapisan bawah
func SetupRoutes(app *fiber.App, db *sql.DB, jwtSecret string) {
	// 1. Buat instance middleware dengan secret key
	mw := middleware.NewMiddleware(jwtSecret)

	// 2. Buat instance handlers dengan dependensi yang sesuai
	authHandler := NewAuthHandler(db, jwtSecret)
	alumniHandler := NewAlumniHandler(db)
	pekerjaanHandler := NewPekerjaanHandler(db)

	// Grup utama API
	api := app.Group("/alumni-management-system")

	// Rute publik untuk login
	api.Post("/login", authHandler.Login)

	// Grup untuk rute yang terproteksi, menggunakan method dari instance middleware
	protected := api.Group("/", mw.JWTMiddleware())
	
	// Rute profil
	protected.Get("/profile", middleware.AdminOrUser(), authHandler.Profile)

	// Rute untuk Alumni
	alumniRoutes := protected.Group("/alumni")
	alumniRoutes.Get("/", middleware.AdminOrUser(), alumniHandler.GetAll)
	alumniRoutes.Get("/:id", middleware.AdminOrUser(), alumniHandler.GetByID)
	alumniRoutes.Get("/laporan/baru-bekerja", middleware.AdminOrUser(), alumniHandler.GetAlumniBaruBekerja)
	alumniRoutes.Post("/", middleware.AdminOnly(), alumniHandler.Create)
	alumniRoutes.Put("/:id", middleware.AdminOnly(), alumniHandler.Update)
	alumniRoutes.Delete("/:id", middleware.AdminOnly(), alumniHandler.Delete)

	// Rute untuk Pekerjaan
	pekerjaanRoutes := protected.Group("/pekerjaan")
	pekerjaanRoutes.Get("/", middleware.AdminOrUser(), pekerjaanHandler.GetAll)
	pekerjaanRoutes.Get("/:id", middleware.AdminOrUser(), pekerjaanHandler.GetByID)
	pekerjaanRoutes.Get("/alumni/:alumni_id", middleware.AdminOnly(), pekerjaanHandler.GetByAlumniID)
	pekerjaanRoutes.Post("/", middleware.AdminOnly(), pekerjaanHandler.Create)
	pekerjaanRoutes.Put("/:id", middleware.AdminOnly(), pekerjaanHandler.Update)
	pekerjaanRoutes.Delete("/:id", middleware.AdminOnly(), pekerjaanHandler.Delete)

	// Rute Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})
}

