package main

import (
	"log"
	"os"
	"pastebin-backend/internal/database"
	"pastebin-backend/internal/router"
	"pastebin-backend/internal/utils/cron"

	"github.com/joho/godotenv"
)

func main() {
	var err error;
    if os.Getenv("RENDER") == "" {
        err = godotenv.Load()
        if err != nil {
            log.Fatal("Error loading .env file")
        }
    }

	database.Connect()
	database.Migrate()

	r := router.SetupRouter()
	cron.StartUploadCleanupJob()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if not provided (for local development)
	}

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
