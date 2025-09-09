package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/entity"
	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/service"
	"github.com/pedroaugustou/qrcode-generator-api/internal/infrastructure/repository"
	"github.com/pedroaugustou/qrcode-generator-api/internal/presentation/dto"
	goqrcode "github.com/skip2/go-qrcode"

	"gorm.io/gorm"
)

type QRCodeUseCase interface {
	GetAllQRCodes(ctx context.Context) ([]dto.QRCodeResponse, error)
	GetQRCodeById(ctx context.Context, id string) (*dto.QRCodeResponse, error)
	AddQRCode(ctx context.Context, req *dto.QRCodeRequest) (*dto.QRCodeResponse, error)
	DeleteQRCode(ctx context.Context, id string) error
}

type qrCodeUseCase struct {
	qrCodeRepository repository.QRCodeRepository
	storageService   service.StorageService
}

func NewQRCodeUseCase(r repository.QRCodeRepository, s service.StorageService) QRCodeUseCase {
	return &qrCodeUseCase{
		qrCodeRepository: r,
		storageService:   s,
	}
}

func (u *qrCodeUseCase) GetAllQRCodes(ctx context.Context) ([]dto.QRCodeResponse, error) {
	qrCodes, err := u.qrCodeRepository.GetAllQRCodes(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.QRCodeResponse, len(qrCodes))
	for i, qrCode := range qrCodes {
		response[i] = *entityToDTO(&qrCode)
	}

	return response, nil
}

func (u *qrCodeUseCase) GetQRCodeById(ctx context.Context, id string) (*dto.QRCodeResponse, error) {
	qrCode, err := u.qrCodeRepository.GetQRCodeById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return entityToDTO(qrCode), nil
}

func (u *qrCodeUseCase) AddQRCode(ctx context.Context, req *dto.QRCodeRequest) (*dto.QRCodeResponse, error) {
	var errs []string

	if req.Content == nil {
		errs = append(errs, "content is required")
	} else if len(*req.Content) < 10 || len(*req.Content) > 255 {
		errs = append(errs, "content must be between 10 and 255 characters")
	}

	if req.RecoveryLevel == nil {
		errs = append(errs, "recovery_level is required")
	} else if *req.RecoveryLevel < 0 || *req.RecoveryLevel > 3 {
		errs = append(errs, "recovery_level must be between 0 and 3")
	}

	if req.Size == nil {
		errs = append(errs, "size is required")
	} else if *req.Size < 256 || *req.Size > 1024 {
		errs = append(errs, "size must be between 256 and 1024")
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("validation failed: %s", strings.Join(errs, "; "))
	}

	qrCode := entity.NewQRCode(*req.Content, *req.Size, *req.RecoveryLevel)

	png, err := goqrcode.Encode(*req.Content, goqrcode.RecoveryLevel(*req.RecoveryLevel), *req.Size)
	if err != nil {
		return nil, err
	}

	url, err := u.storageService.PutQRCode(ctx, png, qrCode)
	if err != nil {
		return nil, err
	}

	qrCode.URL = url

	err = u.qrCodeRepository.AddQRCode(ctx, qrCode)
	if err != nil {
		return nil, err
	}

	return entityToDTO(qrCode), nil
}

func (u *qrCodeUseCase) DeleteQRCode(ctx context.Context, id string) error {
	qrCode, err := u.qrCodeRepository.GetQRCodeById(ctx, id)
	if err != nil {
		return err
	}

	if err := u.storageService.DeleteQRCode(ctx, qrCode.URL); err != nil {
		log.Printf("failed to delete QR code file from storage: %v", err)
	}

	return u.qrCodeRepository.DeleteQRCode(ctx, id)
}

func entityToDTO(e *entity.QRCode) *dto.QRCodeResponse {
	return &dto.QRCodeResponse{
		ID:            e.ID,
		URL:           e.URL,
		Content:       e.Content,
		Size:          e.Size,
		RecoveryLevel: e.RecoveryLevel,
		CreatedAt:     e.CreatedAt,
		ExpiresAt:     e.ExpiresAt,
	}
}
