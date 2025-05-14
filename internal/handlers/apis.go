package handlers

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"pastebin-backend/internal/database"
	"pastebin-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var jwtSecret = []byte("your_secret_key")

// func CreatePaste(c *gin.Context) {
    
//     var newPaste models.Paste
//     if err := c.ShouldBindJSON(&newPaste); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     newPaste.CreatedAt = time.Now()
//     newPaste.UpdatedAt = time.Now()

//     if err := database.DB.Create(&newPaste).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusCreated, newPaste)
// }


// func GetAllPastes(c *gin.Context) {
//     var pastes []models.Paste
//     if err := database.DB.Find(&pastes).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }
//     c.JSON(http.StatusOK, pastes)
// }

// func UpdatePaste(c *gin.Context) {
//     idParam := c.Param("id")
//     id, err := uuid.Parse(idParam)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
//         return
//     }

//     var paste models.Paste
//     if err := database.DB.First(&paste, "id = ?", id).Error; err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Paste not found"})
//         return
//     }

//     var input models.Paste
//     if err := c.ShouldBindJSON(&input); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     paste.Content = input.Content
//     paste.UpdatedAt = time.Now()

//     if err := database.DB.Save(&paste).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, paste)
// }

// func DeletePaste(c *gin.Context) {
//     idParam := c.Param("id")
//     id, err := uuid.Parse(idParam)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
//         return
//     }

//     if err := database.DB.Delete(&models.Paste{}, "id = ?", id).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "Paste deleted"})
// }

func GetPaste(c *gin.Context) {
    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }


    var paste models.Paste
    if err := database.DB.First(&paste, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Paste not found"})
        return
    }

    c.JSON(http.StatusOK, paste)
}

// ----------------------------------------------------------------------------


func CreatePaste(c *gin.Context) {

    fmt.Println("User ID: ",  c.MustGet("userID"))   
    fmt.Println("User ID type: ", reflect.TypeOf( c.MustGet("userID")))

    userIDFloat := c.MustGet("userID").(float64)
    userID := uint(userIDFloat)

    var newPaste models.Paste
    if err := c.ShouldBindJSON(&newPaste); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newPaste.UserID = userID
    newPaste.CreatedAt = time.Now()
    newPaste.UpdatedAt = time.Now()

    if err := database.DB.Create(&newPaste).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, newPaste)
}

func GetAllPastes(c *gin.Context) {
    userIDFloat := c.MustGet("userID").(float64)
    userID := uint(userIDFloat)

    var pastes []models.Paste
    if err := database.DB.Where("user_id = ? order by created_at desc", userID).Find(&pastes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, pastes)
}




func UpdatePaste(c *gin.Context) {
    userIDFloat := c.MustGet("userID").(float64)
    userID := uint(userIDFloat)
    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }

    var paste models.Paste
    if err := database.DB.First(&paste, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Paste not found"})
        return
    }

    // Authorization check
    if paste.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this paste"})
        return
    }

    var input models.Paste
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    paste.Title = input.Title
    paste.Content = input.Content
    paste.UpdatedAt = time.Now()

    if err := database.DB.Save(&paste).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, paste)
}

func DeletePaste(c *gin.Context) {
    userIDFloat := c.MustGet("userID").(float64)
    userID := uint(userIDFloat)
    idParam := c.Param("id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
        return
    }

    var paste models.Paste
    if err := database.DB.First(&paste, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Paste not found"})
        return
    }

    // Authorization check
    if paste.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this paste"})
        return
    }

    if paste.Type == "file" {
        filePath := "./"+paste.Content // Assuming Content holds the file path

       
        // Check if file exists and delete it
        if _, err := os.Stat(filePath); err == nil {
          
            if err := os.Remove(filePath); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete file: %s", err.Error())})
                return
            }
            fmt.Println("File deleted:", filePath)
        } else if !os.IsNotExist(err) {
            c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error checking file: %s", err.Error())})
            return
        }
    }

    if err := database.DB.Delete(&paste).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Paste deleted"})
}


