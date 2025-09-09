package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// load database connection
	dbConn, err := database.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	err = database.AutoMigrate(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	// load storage service
	azureContainer, err := storage.NewAzureBlobConnection()
	if err != nil {
		log.Fatal(err)
	}
	storageService := service.NewStorageService(azureContainer)

	// load /qr
	qrr := repository.NewQRCodeRepository(dbConn)
	qru := usecase.NewQRCodeUseCase(qrr, storageService)
	qrh := handler.NewQRCodeHandler(qru)
	router.SetupQRCodeRoutes(r, qrh)

	// run api
	r.Run(":8080")
}
