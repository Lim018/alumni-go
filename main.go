package main

import (
    "flag"
    "log"
    "os"
    
    "go-fiber/config"
    "go-fiber/database"
    "go-fiber/routes"
    
    _ "go-fiber/docs"
    
    fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Alumni Management API
// @version 1.0
// @description API untuk manajemen data alumni, pekerjaan, dan file
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
    // Command line flags
    migrate := flag.Bool("migrate", false, "Run database migrations")
    seed := flag.Bool("seed", false, "Seed database with initial data")
    reset := flag.Bool("reset", false, "Drop all collections and re-migrate")
    summary := flag.Bool("summary", false, "Show database summary")
    flag.Parse()

    // Load environment variables
    config.LoadEnv()

    // Connect to database
    db := database.ConnectDB()

    // Handle command line operations
    if *reset {
        log.Println("⚠️  Resetting database...")
        if err := database.DropAllCollections(db); err != nil {
            log.Fatal("Failed to drop collections:", err)
        }
        log.Println("✅ Database reset completed")
        return
    }

    if *migrate {
        log.Println("🔄 Running migrations...")
        if err := database.RunMigrations(db); err != nil {
            log.Fatal("Migration failed:", err)
        }
        log.Println("✅ Migrations completed")
        return
    }

    if *seed {
        log.Println("🌱 Seeding database...")
        if err := database.SeedData(db); err != nil {
            log.Fatal("Seeding failed:", err)
        }
        if err := database.SeedSummary(db); err != nil {
            log.Fatal("Failed to show summary:", err)
        }
        return
    }

    if *summary {
        if err := database.SeedSummary(db); err != nil {
            log.Fatal("Failed to show summary:", err)
        }
        return
    }

    // Regular application startup
    log.Println("🚀 Starting application...")

    // Run migrations automatically on startup (optional)
    if os.Getenv("AUTO_MIGRATE") == "true" {
        if err := database.RunMigrations(db); err != nil {
            log.Printf("⚠️  Auto-migration failed: %v", err)
        }
    }

    // Create Fiber app
    app := config.NewApp(db)

    // Swagger route
    app.Get("/swagger/*", fiberSwagger.WrapHandler)

    // Register routes
    routes.RegisterRoutes(app, db)

    // Start server
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "3000"
    }

    log.Printf("🌐 Server running on http://localhost:%s", port)
    log.Printf("📚 Swagger UI: http://localhost:%s/swagger/index.html", port)
    log.Fatal(app.Listen(":" + port))
}