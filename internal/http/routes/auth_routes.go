package routes

import (
	"github.com/DevPulseLab/salat/internal/handlers"
	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	AuthHandler *handlers.AuthHandler
}

func NewAuthRoutes(authHandler *handlers.AuthHandler) *AuthRoutes {
	return &AuthRoutes{AuthHandler: authHandler}
}

func (r *AuthRoutes) Setup(rg *gin.RouterGroup) {
	register := rg.Group("/register")
	{
		register.POST("/cloudflare", r.AuthHandler.CloudflareSSO)
		register.POST("", r.AuthHandler.Register)
	}
	rg.POST("/login", r.AuthHandler.Login)
}
