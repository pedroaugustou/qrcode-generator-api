package dto

import "time"

type QRCodeRequest struct {
	Content       *string `json:"content"`
	RecoveryLevel *int    `json:"recovery_level"`
	Size          *int    `json:"size"`
}

type QRCodeResponse struct {
	ID            string    `json:"id"`
	URL           string    `json:"url"`
	Content       string    `json:"content"`
	Size          int       `json:"size"`
	RecoveryLevel int       `json:"recovery_level"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}
