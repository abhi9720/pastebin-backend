package handlers

import (
	"net/http"
	"pastebin-backend/internal/database"
	"pastebin-backend/internal/models"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFilePaste(c *gin.Context) {
	userIDFloat := c.MustGet("userID").(float64)
	userID := uint(userIDFloat)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file found"})
		return
	}

	title, ok := c.GetPostForm("title")
	if !ok {
		title = file.Filename
	}

	// Save to ./uploads/
	filename := uuid.New().String() + "-" + filepath.Base(file.Filename)
	diskPath := "./uploads/" + filename // for saving to disk
	publicURL := "/uploads/" + filename // for serving via HTTP
	if err := c.SaveUploadedFile(file, diskPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	if title == "" {
		title = file.Filename
	}

	// Save Paste record
	paste := models.Paste{
		Title:     title,
		Content:   publicURL,
		Type:      "file",
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&paste).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create paste record"})
		return
	}

	c.JSON(http.StatusOK, paste)
}
