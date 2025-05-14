package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Paste struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Type      string    `json:"type"`
    UserID    uint      `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (p *Paste) BeforeCreate(tx *gorm.DB) (err error) {
    p.ID = uuid.New()
    if p.Type == "" {
		p.Type = "text"
	}
    return
}
