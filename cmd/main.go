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
    
	envPath := "/etc/secrets/.env"
    if _, err := os.Stat(envPath); os.IsNotExist(err) {
        envPath = ".env"
    }

    // Load environment variables from the .env file
    if err := godotenv.Load(envPath); err != nil {
        log.Printf("No .env file found at %s", envPath)
    }
	
	port := os.Getenv("PORT")
  

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
