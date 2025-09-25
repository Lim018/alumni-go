package main

import (
	"alumni-go/config"
	"alumni-go/database"
	"log"
	// "os"
)

func main() {
	// Muat environment variables
	config.LoadEnv()

	// Buat koneksi database
	db := database.Connect()
	// defer db.Close() // Anda bisa menambahkan ini jika perlu

	// Ambil JWT Secret Key SATU KALI di sini
	jwtSecret := config.GetEnv("JWT_SECRET", "your-super-secret-key-that-is-long")

	// Buat aplikasi Fiber dan inject 'db' dan 'jwtSecret'
	app := config.NewApp(db, jwtSecret)

	// Jalankan server
	port := config.GetEnv("APP_PORT", "3000")
	log.Printf("Server starting on port %s", port)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
