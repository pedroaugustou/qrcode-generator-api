package repository

import (
	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/entity"

	"gorm.io/gorm"
)

type QRCodeRepository interface {
	GetAllQRCodes() ([]entity.QRCode, error)
	GetQRCodeById(id string) (*entity.QRCode, error)
	AddQRCode(qrCode *entity.QRCode) error
}

type qrCodeRepository struct {
	db *gorm.DB
}

func NewQRCodeRepository(db *gorm.DB) QRCodeRepository {
	return &qrCodeRepository{db: db}
}

func (q *qrCodeRepository) GetAllQRCodes() ([]entity.QRCode, error) {
	var qrCodes []entity.QRCode
	result := q.db.Find(&qrCodes)
	return qrCodes, result.Error
}

func (q *qrCodeRepository) GetQRCodeById(id string) (*entity.QRCode, error) {
	var qrCode entity.QRCode
	result := q.db.Where("id = ?", id).First(&qrCode)
	return &qrCode, result.Error
}

func (q *qrCodeRepository) AddQRCode(qrCode *entity.QRCode) error {
	return q.db.Create(qrCode).Error
}
