package entity

import (
	"time"

	"github.com/google/uuid"
)

type QRCode struct {
	ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
	Content       string    `gorm:"not null" json:"content"`
	URL           string    `gorm:"not null" json:"url"`
	Size          int       `gorm:"not null" json:"size"`
	RecoveryLevel int       `gorm:"not null" json:"recovery_level"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

func NewQRCode(content string, size int, recoveryLevel int) *QRCode {
	return &QRCode{
		ID:            uuid.New().String(),
		Content:       content,
		Size:          size,
		RecoveryLevel: recoveryLevel,
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(10 * time.Second),
	}
}
