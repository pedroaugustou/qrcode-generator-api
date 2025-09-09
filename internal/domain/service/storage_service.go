package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/entity"
)

type StorageService interface {
	PutQRCode(ctx context.Context, png []byte, qrCode *entity.QRCode) (string, error)
	DeleteQRCode(ctx context.Context, url string) error
}

type storageService struct {
	containerURL  azblob.ContainerURL
	endpoint      string
	containerName string
}

func NewStorageService(containerURL azblob.ContainerURL) StorageService {
	return &storageService{
		containerURL:  containerURL,
		endpoint:      os.Getenv("AZURE_STORAGE_ENDPOINT"),
		containerName: os.Getenv("AZURE_STORAGE_CONTAINER_NAME"),
	}
}

func (s *storageService) PutQRCode(ctx context.Context, png []byte, qrCode *entity.QRCode) (string, error) {
	blobName := qrCode.ID + ".png"
	blobURL := s.containerURL.NewBlockBlobURL(blobName)

	headers := azblob.BlobHTTPHeaders{
		ContentType: "image/png",
	}

	_, err := azblob.UploadBufferToBlockBlob(
		ctx,
		png,
		blobURL,
		azblob.UploadToBlockBlobOptions{
			BlobHTTPHeaders: headers,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.containerName, blobName), nil
}

func (s *storageService) DeleteQRCode(ctx context.Context, url string) error {
	prefix := fmt.Sprintf("%s/%s/", s.endpoint, s.containerName)

	if !strings.HasPrefix(url, prefix) {
		return fmt.Errorf("invalid URL format: %s", url)
	}

	blobName := strings.TrimPrefix(url, prefix)
	if blobName == "" {
		return fmt.Errorf("blob name is empty")
	}

	blobURL := s.containerURL.NewBlockBlobURL(blobName)
	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		return fmt.Errorf("failed to delete blob %s: %w", blobName, err)
	}

	return nil
}
