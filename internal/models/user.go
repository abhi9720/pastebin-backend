package models

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    GoogleID  string    `json:"google_id"`
    Email     string    `gorm:"uniqueIndex" json:"email"`
    Password  string    `json:"-"`
    Name      string    `json:"name"`
    Picture   string    `json:"picture"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
