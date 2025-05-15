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
	    // Load .env only in local development
    if os.Getenv("RENDER") == "" {
        if err := godotenv.Load(); err != nil {
            log.Fatal("Error loading .env file")
        }
    }

    database.Connect()
    database.Migrate()

    r := router.SetupRouter()
    cron.StartUploadCleanupJob()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default to 8080 for local development
    }

    log.Printf("Server running on port %s", port)

    // Explicitly bind to 0.0.0.0
    if err := r.Run("0.0.0.0:" + port); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
