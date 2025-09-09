package storage

import (
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func NewAzureBlobConnection() (azblob.ContainerURL, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return azblob.ContainerURL{}, fmt.Errorf("failed to create credentials: %w", err)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	serviceURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	containerURL := azblob.NewServiceURL(*serviceURL, pipeline).NewContainerURL(containerName)

	return containerURL, nil
}
