package storage

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOConnection() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_ACCESS_SECRET")

	useSSL := false
	if sslEnv := os.Getenv("MINIO_USE_SSL"); sslEnv != "" {
		ssl, err := strconv.ParseBool(sslEnv)
		if err != nil {
			log.Printf("Invalid MINIO_USE_SSL value, defaulting to false: %v", err)
		} else {
			useSSL = ssl
		}
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return client, nil
}
