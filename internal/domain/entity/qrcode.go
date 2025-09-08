package entity

import (
	"time"

	"github.com/google/uuid"
)

type QRCode struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Content   string    `gorm:"not null" json:"content"`
	URL       string    `gorm:"not null" json:"url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

func NewQRCode(content string) *QRCode {
	return &QRCode{
		ID:        uuid.New().String(),
		Content:   content,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	}
}
