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
    if os.Getenv("RENDER") == "" {
        if err := godotenv.Load(); err != nil {
            log.Fatal("Error loading .env file")
        }
    }



	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

	log.Printf("PORT environment variable: %s", port)

    database.Connect()
    database.Migrate()

    r := router.SetupRouter()
    cron.StartUploadCleanupJob()

    log.Printf("Server running on port %s", port)

    // Explicitly bind to 0.0.0.0
    if err := r.Run("0.0.0.0:" + port); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
