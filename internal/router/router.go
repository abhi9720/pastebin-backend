package router

import (
	"os"
	"pastebin-backend/internal/handlers"
	"pastebin-backend/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {    

    r := gin.Default()

     r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{os.Getenv("FrontendURL")},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

     // OAuth routes
    r.GET("/auth/google/login", handlers.GoogleLogin)
    r.GET("/auth/google/callback", handlers.GoogleCallback)

    // Auth routes
    r.POST("/auth/login", handlers.EmailLogin)
    r.POST("/auth/signup", handlers.Register)

    // Upload Routes

    r.Static("/uploads", "./uploads") // serve uploaded files
    r.POST("/api/paste/upload", middleware.AuthMiddleware(), handlers.UploadFilePaste)


    // Define routes
    r.POST("/pastes", middleware.AuthMiddleware(), handlers.CreatePaste)
    r.GET("/pastes/:id",  middleware.AuthMiddleware(), handlers.GetPaste)
    r.GET("/pastes", middleware.AuthMiddleware(), handlers.GetAllPastes)
    r.PUT("/pastes/:id",  middleware.AuthMiddleware(), handlers.UpdatePaste)
    r.DELETE("/pastes/:id",  middleware.AuthMiddleware(), handlers.DeletePaste)

    return r
}