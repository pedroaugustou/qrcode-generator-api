package usecase

import (
	"errors"
	"fmt"
	"qrcode-generator-api/internal/domain/entity"
	"qrcode-generator-api/internal/domain/service"
	"qrcode-generator-api/internal/infrastructure/repository"
	"qrcode-generator-api/internal/presentation/dto"
	"strings"

	goqrcode "github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type QRCodeUseCase interface {
	GetAllQRCodes() ([]dto.QRCodeResponse, error)
	GetQRCodeById(id string) (*entity.QRCode, error)
	AddQRCode(req *dto.CreateQRCodeRequest) (*dto.QRCodeResponse, error)
}

type qrcodeUseCase struct {
	r repository.QRCodeRepository
	s service.StorageService
}

func NewQRCodeUseCase(r repository.QRCodeRepository, s service.StorageService) QRCodeUseCase {
	return &qrcodeUseCase{
		r: r,
		s: s,
	}
}

func (q *qrcodeUseCase) GetAllQRCodes() ([]dto.QRCodeResponse, error) {
	qrcodes, err := q.r.GetAllQRCodes()
	if err != nil {
		return nil, err
	}

	response := make([]dto.QRCodeResponse, len(qrcodes))
	for i, qrcode := range qrcodes {
		response[i] = dto.QRCodeResponse{
			ID:        qrcode.ID,
			URL:       qrcode.URL,
			Content:   qrcode.Content,
			CreatedAt: qrcode.CreatedAt,
			ExpiresAt: qrcode.ExpiresAt,
		}
	}

	return response, nil
}

func (u *qrcodeUseCase) GetQRCodeById(id string) (*entity.QRCode, error) {
	qrcode, err := u.r.GetQRCodeById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return qrcode, nil
}

func (q *qrcodeUseCase) AddQRCode(req *dto.CreateQRCodeRequest) (*dto.QRCodeResponse, error) {
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

	qrcode := entity.NewQRCode(*req.Content)

	png, err := goqrcode.Encode(*req.Content, goqrcode.RecoveryLevel(*req.RecoveryLevel), *req.Size)
	if err != nil {
		return nil, err
	}

	url, err := q.s.PutQRCode(png, qrcode)
	if err != nil {
		return nil, err
	}

	qrcode.URL = url

	err = q.r.AddQRCode(qrcode)
	if err != nil {
		return nil, err
	}

	return &dto.QRCodeResponse{
		ID:        qrcode.ID,
		URL:       qrcode.URL,
		Content:   qrcode.Content,
		CreatedAt: qrcode.CreatedAt,
		ExpiresAt: qrcode.ExpiresAt,
	}, nil
}
