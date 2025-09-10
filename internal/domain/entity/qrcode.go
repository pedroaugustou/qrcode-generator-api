package entity

import (
	"time"

	"github.com/google/uuid"
)

type QRCode struct {
	ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
	Content       string    `json:"content"`
	URL           string    `json:"url"`
	Size          int       `json:"size"`
	RecoveryLevel int       `json:"recovery_level"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewQRCode(content string, size int, recoveryLevel int) *QRCode {
	return &QRCode{
		ID:            uuid.New().String(),
		Content:       content,
		Size:          size,
		RecoveryLevel: recoveryLevel,
		CreatedAt:     time.Now().UTC(),
	}
}
