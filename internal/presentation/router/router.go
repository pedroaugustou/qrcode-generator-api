package router

import (
	"qrcode-generator-api/internal/presentation/handler"

	"github.com/gin-gonic/gin"
)

func SetupQRCodeRoutes(r *gin.Engine, h *handler.QRCodeHandler) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": "OK",
			"data": "pong",
		})
	})

	qr := v1.Group("/qr")
	{
		qr.GET("/:id", h.GetQRCodeById)
		qr.GET("", h.GetAllQRCodes)
		qr.POST("", h.AddQRCode)
	}
}
