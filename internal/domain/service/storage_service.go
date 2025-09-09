package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/entity"
)

type StorageService interface {
	PutQRCode(ctx context.Context, png []byte, qrCode *entity.QRCode) (string, error)
	DeleteQRCode(ctx context.Context, url string) error
	CleanupExpiredFiles(ctx context.Context) error
}

type storageService struct {
	containerURL azblob.ContainerURL
}

func NewStorageService(containerURL azblob.ContainerURL) StorageService {
	return &storageService{containerURL: containerURL}
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

	endpoint := os.Getenv("AZURE_STORAGE_ENDPOINT")
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	return fmt.Sprintf("%s/%s/%s", endpoint, containerName, blobName), nil
}

func (s *storageService) CleanupExpiredFiles(ctx context.Context) error {
	now := time.Now().UTC().Truncate(time.Hour)

	marker := azblob.Marker{}
	for marker.NotDone() {
		listBlob, err := s.containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			return fmt.Errorf("failed to list blobs: %w", err)
		}

		for _, blobInfo := range listBlob.Segment.BlobItems {
			if blobInfo.Properties.LastModified.Before(now.Add(-24 * time.Hour)) {
				blobURL := s.containerURL.NewBlockBlobURL(blobInfo.Name)
				_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
				if err != nil {
					log.Printf("failed to delete blob %s: %v", blobInfo.Name, err)
					continue
				}
				log.Printf("deleted expired blob: %s", blobInfo.Name)
			}
		}

		marker = listBlob.NextMarker
	}

	return nil
}

func (s *storageService) DeleteQRCode(ctx context.Context, url string) error {
	endpoint := os.Getenv("AZURE_STORAGE_ENDPOINT")
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	prefix := fmt.Sprintf("%s/%s/", endpoint, containerName)

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
