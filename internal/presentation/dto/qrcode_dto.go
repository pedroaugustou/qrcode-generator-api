package dto

import "time"

type CreateQRCodeRequest struct {
	Content       *string `json:"content"`
	RecoveryLevel *int    `json:"recovery_level"`
	Size          *int    `json:"size"`
}

type QRCodeResponse struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
