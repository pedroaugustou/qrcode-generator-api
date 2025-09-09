package repository

import (
	"context"
	"time"

	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/entity"

	"gorm.io/gorm"
)

type QRCodeRepository interface {
	GetAllQRCodes(ctx context.Context) ([]entity.QRCode, error)
	GetQRCodeById(ctx context.Context, id string) (*entity.QRCode, error)
	AddQRCode(ctx context.Context, qrCode *entity.QRCode) error
	DeleteQRCode(ctx context.Context, id string) error
	DeleteExpiredQRCodes(ctx context.Context) error
}

type qrCodeRepository struct {
	database *gorm.DB
}

func NewQRCodeRepository(d *gorm.DB) QRCodeRepository {
	return &qrCodeRepository{database: d}
}

func (q *qrCodeRepository) GetAllQRCodes(ctx context.Context) ([]entity.QRCode, error) {
	var qrCodes []entity.QRCode
	result := q.database.WithContext(ctx).Find(&qrCodes)
	return qrCodes, result.Error
}

func (q *qrCodeRepository) GetQRCodeById(ctx context.Context, id string) (*entity.QRCode, error) {
	var qrCode entity.QRCode
	result := q.database.WithContext(ctx).Where("id = ?", id).First(&qrCode)
	return &qrCode, result.Error
}

func (q *qrCodeRepository) AddQRCode(ctx context.Context, qrCode *entity.QRCode) error {
	return q.database.WithContext(ctx).Create(qrCode).Error
}

func (q *qrCodeRepository) DeleteQRCode(ctx context.Context, id string) error {
	result := q.database.WithContext(ctx).Where("id = ?", id).Delete(&entity.QRCode{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (q *qrCodeRepository) DeleteExpiredQRCodes(ctx context.Context) error {
	now := time.Now().UTC().Truncate(time.Hour)
	return q.database.WithContext(ctx).Where("expires_at <= ?", now).Delete(&entity.QRCode{}).Error
}
