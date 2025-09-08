package repository

import (
	"qrcode-generator-api/internal/domain/entity"

	"gorm.io/gorm"
)

type QRCodeRepository interface {
	GetAllQRCodes() ([]entity.QRCode, error)
	GetQRCodeById(id string) (*entity.QRCode, error)
	AddQRCode(qrcode *entity.QRCode) error
}

type qrcodeRepository struct {
	db *gorm.DB
}

func NewQRCodeRepository(db *gorm.DB) QRCodeRepository {
	return &qrcodeRepository{db: db}
}

func (q *qrcodeRepository) GetAllQRCodes() ([]entity.QRCode, error) {
	var qrcodes []entity.QRCode
	result := q.db.Find(&qrcodes)
	return qrcodes, result.Error
}

func (q *qrcodeRepository) GetQRCodeById(id string) (*entity.QRCode, error) {
	var qrcode entity.QRCode
	result := q.db.Where("id = ?", id).First(&qrcode)
	return &qrcode, result.Error
}

func (q *qrcodeRepository) AddQRCode(qrcode *entity.QRCode) error {
	return q.db.Create(qrcode).Error
}
