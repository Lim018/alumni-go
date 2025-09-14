package main

import (
	"alumni-management-system/config"
	"alumni-management-system/database"
	"log"
)

func main() {
	config.LoadEnv()

	database.Connect()
	defer database.DB.Close()

	app := config.CreateApp()

	port := config.GetEnv("APP_PORT", "3000")
	log.Printf("Server starting on port %s", port)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}