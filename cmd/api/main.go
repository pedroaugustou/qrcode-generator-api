package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/service"
	"github.com/pedroaugustou/qrcode-generator-api/internal/infrastructure/database"
	"github.com/pedroaugustou/qrcode-generator-api/internal/infrastructure/repository"
	"github.com/pedroaugustou/qrcode-generator-api/internal/infrastructure/storage"
	"github.com/pedroaugustou/qrcode-generator-api/internal/presentation/handler"
	"github.com/pedroaugustou/qrcode-generator-api/internal/presentation/router"
	"github.com/pedroaugustou/qrcode-generator-api/internal/usecase"
)

func main() {
	r := gin.Default()

	// load database
	dbConn, err := database.NewDBConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := database.AutoMigrate(dbConn); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// load storage service
	azureContainer, err := storage.NewAzureBlobConnection()
	if err != nil {
		log.Fatalf("failed to connect to azure storage: %v", err)
	}
	storageService := service.NewStorageService(azureContainer)

	// load repos, handlers, usecases
	qrr := repository.NewQRCodeRepository(dbConn)
	qru := usecase.NewQRCodeUseCase(qrr, storageService)
	qrh := handler.NewQRCodeHandler(qru)

	// setup api routes
	router.SetupQRCodeRoutes(r, qrh)

	r.Run(":8080")
}
