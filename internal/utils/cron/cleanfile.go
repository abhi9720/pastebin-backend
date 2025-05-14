package cron

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func StartUploadCleanupJob() {
	ticker := time.NewTicker(5 * time.Minute) // runs every 5 Minutes
	go func() {
		log.Println("Starting cleanup job...")
		for {
			<-ticker.C
			log.Println("Running cleanup task...")
			cleanupUploads("./uploads")
		}
	}()
}

func cleanupUploads(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Println("Error reading upload dir:", err)
		return
	}

	now := time.Now()
	for _, file := range files {
		fmt.Print(file.Name())
		path := filepath.Join(dir, file.Name())
		fmt.Println(path)
		info, err := os.Stat(path)

		if err != nil {
			continue
		}
		if now.Sub(info.ModTime()) > 6*time.Hour {
			if err := os.Remove(path); err != nil {
				log.Println("Error deleting file:", err)
			} else {
				log.Println("Deleted old file:", path)
			}
		}
	}
}
