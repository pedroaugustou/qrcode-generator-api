package service

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/entity"
)

type StorageService interface {
	PutQRCode(png []byte, qrCode *entity.QRCode) (string, error)
}

type storageService struct {
	containerURL azblob.ContainerURL
}

func NewStorageService(containerURL azblob.ContainerURL) StorageService {
	return &storageService{containerURL: containerURL}
}

func (s *storageService) PutQRCode(png []byte, qrCode *entity.QRCode) (string, error) {
	blobName := qrCode.ID + ".png"
	blobURL := s.containerURL.NewBlockBlobURL(blobName)

	headers := azblob.BlobHTTPHeaders{
		ContentType: "image/png",
	}

	_, err := azblob.UploadBufferToBlockBlob(
		context.Background(),
		png,
		blobURL,
		azblob.UploadToBlockBlobOptions{
			BlobHTTPHeaders: headers,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	endpoint := os.Getenv("AZURE_STORAGE_ENDPOINT")
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	return fmt.Sprintf("%s/%s/%s", endpoint, containerName, blobName), nil
}
