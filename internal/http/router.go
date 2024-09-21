package http

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializeRoutes(router *gin.Engine, db *gorm.DB, config *config.Config) {
	authHandler := handlers.NewAuthHandler(db, config)

	router.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"ping": "pong"})
	})

	router.POST("/api/register", authHandler.Register)
	router.POST("/api/login", authHandler.Login)
}
