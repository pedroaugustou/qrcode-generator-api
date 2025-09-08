package main

import (
	"log"
	"qrcode-generator-api/internal/domain/service"
	"qrcode-generator-api/internal/infrastructure/database"
	"qrcode-generator-api/internal/infrastructure/repository"
	"qrcode-generator-api/internal/infrastructure/storage"
	"qrcode-generator-api/internal/presentation/handler"
	"qrcode-generator-api/internal/presentation/router"
	"qrcode-generator-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	// load environment variables
	err := godotenv.Load("../../.env")
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
	minioConn, err := storage.NewMinIOConnection()
	if err != nil {
		log.Fatal(err)
	}
	storageService := service.NewStorageService(minioConn)

	// load /qr
	qrr := repository.NewQRCodeRepository(dbConn)
	qru := usecase.NewQRCodeUseCase(qrr, storageService)
	qrh := handler.NewQRCodeHandler(qru)
	router.SetupQRCodeRoutes(r, qrh)

	// run api
	r.Run(":8080")
}
