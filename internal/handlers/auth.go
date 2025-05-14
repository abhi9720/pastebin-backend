package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"pastebin-backend/internal/auth"
	"pastebin-backend/internal/database"
	"pastebin-backend/internal/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name"`  // Added Name field
}

type AuthResponse struct {
    Token string       `json:"token"`
    User  models.User  `json:"user"`
}

func Register(c *gin.Context) {
    var input AuthInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    input.Email = strings.ToLower(input.Email)

    // Check if user already exists
    var existing models.User
    result := database.DB.Where("email = ?", input.Email).First(&existing)
    if result.Error == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user := models.User{
        Email:    input.Email,
        Password: string(hashedPassword),
        Name:     input.Name,  // Added Name to user model
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Return JWT
    token, err := generateJWT(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    user.Password = ""  // Don't return password in the response
    c.JSON(http.StatusCreated, AuthResponse{Token: token, User: user})
}

func EmailLogin(c *gin.Context) {
    var input AuthInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    input.Email = strings.ToLower(input.Email)

    var user models.User
    if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := generateJWT(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    user.Password = ""  // Don't return password in the response
    c.JSON(http.StatusOK, AuthResponse{Token: token, User: user})
}



func GoogleLogin(c *gin.Context) {
    url := auth.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
        return
    }

    ctx := context.Background()

    // Exchange code for tokens
    token, err := auth.GoogleOAuthConfig.Exchange(ctx, code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
        return
    }

    idToken := token.Extra("id_token")
    idTokenStr, ok := idToken.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ID token"})
        return
    }
   

    // Verify ID token manually
    payload, err := verifyIDToken(idTokenStr)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
        return
    }

    // Extract user info from payload
    userEmail := payload["email"].(string)
    userName := ""
    if name, ok := payload["name"].(string); ok {
        userName = name
    }
    userID := payload["sub"].(string)
    picture := ""
    if p, ok := payload["picture"].(string); ok {
        picture = p
    }


    // Check if user exists
    var user models.User
    result := database.DB.Where("google_id = ?", userID).First(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            // Create new user
            user = models.User{
                GoogleID: userID,
                Email:    userEmail,
                Name:     userName,
                Picture:  picture,
            }
            if err := database.DB.Create(&user).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
                return
            }
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            return
        }
    } else {
        // Optionally, update user info if changed
        user.Email = userEmail
        user.Name = userName
        user.Picture = picture
        database.DB.Save(&user)
    }


    // Generate JWT for your app
    tokenString, err := generateJWT(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    origin := c.Request.Header.Get("Origin")
    if origin == "" {
        origin = os.Getenv("FrontendURL")
    }

    // Redirect with the token to the same domain
    frontendURL := fmt.Sprintf("%s/auth/google/callback?token=%s", origin, tokenString)
    c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

func generateJWT(user models.User) (string, error){
    claims := jwt.MapClaims{
        "sub":   user.ID,
        "email": user.Email,
        "name":  user.Name,
        "picture": user.Picture, 
        "exp":   time.Now().Add(72 * time.Hour).Unix(), // token expiry
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

func verifyIDToken(idToken string) (jwt.MapClaims, error) {
    resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
       return nil, errors.New("invalid ID token")
    }

    var payload map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
        return nil, err
    }

    // Cast to jwt.MapClaims
    claims := jwt.MapClaims(payload)
    return claims, nil
}