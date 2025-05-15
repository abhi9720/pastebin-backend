package auth

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func init() {
	var err error;
    if os.Getenv("RENDER") == "" {
        err = godotenv.Load()
        if err != nil {
            log.Fatal("Error loading .env file")
        }
    }
    
    GoogleOAuthConfig = &oauth2.Config{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("BackendURL")+"/auth/google/callback",
        Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
        Endpoint:     google.Endpoint,
    }
}