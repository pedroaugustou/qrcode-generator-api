package dto

import (
	"time"

	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/entity"
)

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
}

func (q *QRCodeResponse) FromEntity(e *entity.QRCode) *QRCodeResponse {
	q.ID = e.ID
	q.URL = e.URL
	q.Content = e.Content
	q.Size = e.Size
	q.RecoveryLevel = e.RecoveryLevel
	q.CreatedAt = e.CreatedAt
	return q
}
