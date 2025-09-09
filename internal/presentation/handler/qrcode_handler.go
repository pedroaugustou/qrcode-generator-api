package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pedroaugustou/qrcode-generator-api/internal/presentation/dto"
	"github.com/pedroaugustou/qrcode-generator-api/internal/usecase"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type QRCodeHandler struct {
	useCase usecase.QRCodeUseCase
}

func NewQRCodeHandler(useCase usecase.QRCodeUseCase) *QRCodeHandler {
	return &QRCodeHandler{useCase: useCase}
}

func (h *QRCodeHandler) GetAllQRCodes(ctx *gin.Context) {
	qrCodes, err := h.useCase.GetAllQRCodes(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": qrCodes})
}

func (h *QRCodeHandler) GetQRCodeById(ctx *gin.Context) {
	qrCode, err := h.useCase.GetQRCodeById(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if qrCode == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "qr code not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": qrCode})
}

func (h *QRCodeHandler) AddQRCode(ctx *gin.Context) {
	var req dto.QRCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
		return
	}

	response, err := h.useCase.AddQRCode(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": response})
}

func (h *QRCodeHandler) DeleteQRCode(ctx *gin.Context) {
	err := h.useCase.DeleteQRCode(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "qr code not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
