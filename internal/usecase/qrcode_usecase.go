package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/entity"
	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/service"
	"github.com/pedroaugustou/qrcode-generator-api/internal/infrastructure/repository"
	"github.com/pedroaugustou/qrcode-generator-api/internal/presentation/dto"
	goqrcode "github.com/skip2/go-qrcode"

	"gorm.io/gorm"
)

type QRCodeUseCase interface {
	GetAllQRCodes() ([]dto.QRCodeResponse, error)
	GetQRCodeById(id string) (*entity.QRCode, error)
	AddQRCode(req *dto.CreateQRCodeRequest) (*dto.QRCodeResponse, error)
}

type qrCodeUseCase struct {
	r repository.QRCodeRepository
	s service.StorageService
}

func NewQRCodeUseCase(r repository.QRCodeRepository, s service.StorageService) QRCodeUseCase {
	return &qrCodeUseCase{
		r: r,
		s: s,
	}
}

func (q *qrCodeUseCase) GetAllQRCodes() ([]dto.QRCodeResponse, error) {
	qrCodes, err := q.r.GetAllQRCodes()
	if err != nil {
		return nil, err
	}

	response := make([]dto.QRCodeResponse, len(qrCodes))
	for i, qrCode := range qrCodes {
		response[i] = dto.QRCodeResponse{
			ID:        qrCode.ID,
			URL:       qrCode.URL,
			Content:   qrCode.Content,
			CreatedAt: qrCode.CreatedAt,
			ExpiresAt: qrCode.ExpiresAt,
		}
	}

	return response, nil
}

func (u *qrCodeUseCase) GetQRCodeById(id string) (*entity.QRCode, error) {
	qrCode, err := u.r.GetQRCodeById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return qrCode, nil
}

func (q *qrCodeUseCase) AddQRCode(req *dto.CreateQRCodeRequest) (*dto.QRCodeResponse, error) {
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

	qrCode := entity.NewQRCode(*req.Content)

	png, err := goqrcode.Encode(*req.Content, goqrcode.RecoveryLevel(*req.RecoveryLevel), *req.Size)
	if err != nil {
		return nil, err
	}

	url, err := q.s.PutQRCode(png, qrCode)
	if err != nil {
		return nil, err
	}

	qrCode.URL = url

	err = q.r.AddQRCode(qrCode)
	if err != nil {
		return nil, err
	}

	return &dto.QRCodeResponse{
		ID:        qrCode.ID,
		URL:       qrCode.URL,
		Content:   qrCode.Content,
		CreatedAt: qrCode.CreatedAt,
		ExpiresAt: qrCode.ExpiresAt,
	}, nil
}
