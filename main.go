package main

import (
	"alumni-go/app"
	"alumni-go/config"
	"alumni-go/database"
	"log"
)

func main() {
	config.LoadEnv()

	database.Connect()
	defer database.DB.Close()

	webApp := app.CreateApp()

	port := config.GetEnv("APP_PORT", "3000")
	log.Printf("Server starting on port %s", port)
	
	if err := webApp.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}