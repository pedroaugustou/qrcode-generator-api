package repository

import (
	"qrcode-generator-api/internal/domain/entity"

	"gorm.io/gorm"
)

type QRCodeRepository interface {
	AddQRCode(qrcode *entity.QRCode) error
}

type qrcodeRepository struct {
	db *gorm.DB
}

func NewQRCodeRepository(db *gorm.DB) QRCodeRepository {
	return &qrcodeRepository{db: db}
}

func (q *qrcodeRepository) AddQRCode(qrcode *entity.QRCode) error {
	return q.db.Create(qrcode).Error
}
