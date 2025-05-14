package database

import (
	"log"
	"os"
	"pastebin-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	 var err error

    // Get the connection string from environment variable
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL environment variable not set")
    }

   
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }
}

func Migrate() {
    err := DB.AutoMigrate(&models.Paste{})
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }
    
    err = DB.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }
}