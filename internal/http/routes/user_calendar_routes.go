package routes

import (
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/handlers"
	"github.com/DevPulseLab/salat/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type UserCalendarRoutes struct {
	UserCalendarHandler   *handlers.UserCalendarHandler
	JwtMiddleware         *middlewares.JwtMiddleware
	RoleMiddleware        *middlewares.RoleMiddleware
	CurrentUserMiddleware *middlewares.CurrentUserMiddleware
}

func NewUserCalendarRoutes(
	userCalendarHandler *handlers.UserCalendarHandler,
	jwtMiddleware *middlewares.JwtMiddleware,
	roleMiddleware *middlewares.RoleMiddleware,
	currentUserMiddleware *middlewares.CurrentUserMiddleware,
) *UserCalendarRoutes {
	return &UserCalendarRoutes{
		UserCalendarHandler:   userCalendarHandler,
		JwtMiddleware:         jwtMiddleware,
		RoleMiddleware:        roleMiddleware,
		CurrentUserMiddleware: currentUserMiddleware,
	}
}

func (r *UserCalendarRoutes) Setup(rg *gin.RouterGroup) {
	userCalendar := rg.Group("/user/calendar")
	userCalendar.Use(r.JwtMiddleware.Process)
	{
		userCalendar.POST("/add", r.RoleMiddleware.Process(models.RoleUser), r.CurrentUserMiddleware.Process, r.UserCalendarHandler.Add)
		userCalendar.GET("/current-user-list", r.RoleMiddleware.Process(models.RoleUser), r.UserCalendarHandler.CurrentUserList)
		userCalendar.POST("/remove-for-current-user", r.RoleMiddleware.Process(models.RoleUser), r.UserCalendarHandler.RemoveEntryForCurrentUser)
		userCalendar.GET("/get-close-intervals", r.RoleMiddleware.Process(models.RoleUser, models.RoleAdmin), r.UserCalendarHandler.GetCloseDateInterval)
	}
}
