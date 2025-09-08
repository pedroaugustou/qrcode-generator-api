package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"qrcode-generator-api/internal/domain/entity"

	"github.com/minio/minio-go/v7"
)

type StorageService interface {
	PutQRCode(png []byte, qrcode *entity.QRCode) (string, error)
}

type storageService struct {
	m *minio.Client
}

func NewStorageService(m *minio.Client) StorageService {
	return &storageService{m: m}
}

func (s *storageService) PutQRCode(png []byte, qrcode *entity.QRCode) (string, error) {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	objectName := qrcode.ID + ".png"
	reader := bytes.NewReader(png)

	opts := minio.PutObjectOptions{ContentType: "image/png", Expires: qrcode.ExpiresAt}
	_, err := s.m.PutObject(context.Background(), bucketName, objectName, reader, int64(len(png)), opts)
	if err != nil {
		return "", err
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	return fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, objectName), nil
}
