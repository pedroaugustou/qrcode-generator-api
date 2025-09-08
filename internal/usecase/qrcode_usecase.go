package usecase

import (
	"fmt"
	"qrcode-generator-api/internal/domain/entity"
	"qrcode-generator-api/internal/domain/service"
	"qrcode-generator-api/internal/infrastructure/repository"
	"qrcode-generator-api/internal/presentation/dto"
	"strings"

	goqrcode "github.com/skip2/go-qrcode"
)

type QRCodeUseCase interface {
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

func (q *qrcodeUseCase) AddQRCode(req *dto.CreateQRCodeRequest) (*dto.QRCodeResponse, error) {
	var errors []string

	if req.Content == nil {
		errors = append(errors, "content is required")
	} else if len(*req.Content) < 10 || len(*req.Content) > 255 {
		errors = append(errors, "content must be between 10 and 255 characters")
	}

	if req.RecoveryLevel == nil {
		errors = append(errors, "recovery_level is required")
	} else if *req.RecoveryLevel < 0 || *req.RecoveryLevel > 3 {
		errors = append(errors, "recovery_level must be between 0 and 3")
	}

	if req.Size == nil {
		errors = append(errors, "size is required")
	} else if *req.Size < 256 || *req.Size > 1024 {
		errors = append(errors, "size must be between 256 and 1024")
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
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
		ExpiresAt: &qrcode.ExpiresAt,
	}, nil
}
