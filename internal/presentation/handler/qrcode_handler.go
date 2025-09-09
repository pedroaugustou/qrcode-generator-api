package handler

import (
	"fmt"
	"net/http"

	"github.com/pedroaugustou/qrcode-generator-api/internal/presentation/dto"
	"github.com/pedroaugustou/qrcode-generator-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type QRCodeHandler struct {
	useCase usecase.QRCodeUseCase
}

func NewQRCodeHandler(useCase usecase.QRCodeUseCase) *QRCodeHandler {
	return &QRCodeHandler{useCase: useCase}
}

func (h *QRCodeHandler) GetAllQRCodes(c *gin.Context) {
	qrCodes, err := h.useCase.GetAllQRCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": qrCodes})
}

func (h *QRCodeHandler) GetQRCodeById(c *gin.Context) {
	qrCode, err := h.useCase.GetQRCodeById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if qrCode == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "qr code not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": qrCode})
}

func (h *QRCodeHandler) AddQRCode(c *gin.Context) {
	var req dto.CreateQRCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
		return
	}

	r, err := h.useCase.AddQRCode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": r})
}
