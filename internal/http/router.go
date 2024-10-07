package http

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/handlers"
	"github.com/DevPulseLab/salat/internal/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializeRoutes(router *gin.Engine, db *gorm.DB, config *config.Config) {
	jwtMiddleware := middlewares.NewJwtMiddleware(config)
	roleMiddleware := middlewares.NewRoleMiddleware()
	authHandler := handlers.NewAuthHandler(db, config)
	userHandler := handlers.NewUserHandler(db)

	router.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"ping": "pong"})
	})

	router.POST("/api/register", authHandler.Register)
	router.POST("/api/login", authHandler.Login)

	router.GET("/api/users/list", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), userHandler.GetUserList)
}
