package routes

import (
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/handlers"
	"github.com/DevPulseLab/salat/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	UserHandler           *handlers.UserHandler
	RoleMiddleware        *middlewares.RoleMiddleware
	CurrentUserMiddleware *middlewares.CurrentUserMiddleware
}

func NewUserRoutes(
	userHandler *handlers.UserHandler,
	roleMiddleware *middlewares.RoleMiddleware,
	currentUserMiddleware *middlewares.CurrentUserMiddleware,
) *UserRoutes {
	return &UserRoutes{
		UserHandler:           userHandler,
		RoleMiddleware:        roleMiddleware,
		CurrentUserMiddleware: currentUserMiddleware,
	}
}

func (r *UserRoutes) Setup(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("/me", r.RoleMiddleware.Process(models.RoleUser), r.CurrentUserMiddleware.Process, r.UserHandler.GetCurrentUserInfo)
		users.GET("/list", r.RoleMiddleware.Process(models.RoleAdmin), r.UserHandler.GetUserList)
		users.POST("/set-penalty-card", r.RoleMiddleware.Process(models.RoleAdmin), r.UserHandler.SetPenaltyCard)
	}
}
