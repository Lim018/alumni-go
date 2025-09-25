package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupLogger bisa dipanggil dari NewApp jika diperlukan kustomisasi
// Untuk saat ini, kita akan menggunakan logger bawaan Fiber di NewApp.
// File ini dibuat agar sesuai struktur modul.
func SetupLogger(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
}