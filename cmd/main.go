package main

import (
	"log"
	"os"
	"pastebin-backend/internal/database"
	"pastebin-backend/internal/router"
	"pastebin-backend/internal/utils/cron"
)

func main() {
	
	port := os.Getenv("PORT")
	log.Printf("PORT environment variable: %s", port)

    database.Connect()
    database.Migrate()

    log.Println("Database connected and migrated")

    r := router.SetupRouter()
    cron.StartUploadCleanupJob()

    log.Println("Starting server...")

    // Explicitly bind to 0.0.0.0
    if err := r.Run("0.0.0.0:" + port); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
