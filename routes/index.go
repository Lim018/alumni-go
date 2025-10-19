package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	AlumniRoutes(app, db)
	PekerjaanRoutes(app, db)
	AuthRoutes(app, db) 
	UserRoutes(app, db)
}
