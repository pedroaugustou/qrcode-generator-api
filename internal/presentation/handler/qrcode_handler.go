package handler

import (
	"fmt"
	"net/http"
	"qrcode-generator-api/internal/presentation/dto"
	"qrcode-generator-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type QRCodeHandler struct {
	useCase usecase.QRCodeUseCase
}

func NewQRCodeHandler(useCase usecase.QRCodeUseCase) *QRCodeHandler {
	return &QRCodeHandler{useCase: useCase}
}

func (h *QRCodeHandler) AddQRCode(ctx *gin.Context) {
	var req dto.CreateQRCodeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
		return
	}

	response, err := h.useCase.AddQRCode(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
