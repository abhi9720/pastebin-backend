package main

import (
	"log"
	"pastebin-backend/internal/database"
	"pastebin-backend/internal/router"
	"pastebin-backend/internal/utils/cron"

	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }


    database.Connect()
    database.Migrate()

    r := router.SetupRouter()
    cron.StartUploadCleanupJob()

    // Start server
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}

// expernal postgresql://pastebindb_pncl_user:Psbw8xYEQJfYrmt6X4nrw7PsKF6iBIJI@dpg-d0hgtjh5pdvs73egf4o0-a.oregon-postgres.render.com/pastebindb_pncl