package http

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializeRoutes(router *gin.Engine, db *gorm.DB) {
	authHandler := handlers.NewAuthHandler(db)

	router.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"ping": "pong"})
	})

	router.POST("/api/register", authHandler.Register)
}
