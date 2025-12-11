package routes

import (
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/handlers"
	"github.com/DevPulseLab/salat/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type AdminCalendarRoutes struct {
	AdminCalendarHandler *handlers.AdminCalendarHandler
	JwtMiddleware        *middlewares.JwtMiddleware
	RoleMiddleware       *middlewares.RoleMiddleware
}

func NewAdminCalendarRoutes(
	adminCalendarHandler *handlers.AdminCalendarHandler,
	jwtMiddleware *middlewares.JwtMiddleware,
	roleMiddleware *middlewares.RoleMiddleware,
) *AdminCalendarRoutes {
	return &AdminCalendarRoutes{
		AdminCalendarHandler: adminCalendarHandler,
		JwtMiddleware:        jwtMiddleware,
		RoleMiddleware:       roleMiddleware,
	}
}

func (r *AdminCalendarRoutes) Setup(rg *gin.RouterGroup) {
	adminCalendar := rg.Group("/admin/calendar")
	adminCalendar.Use(r.JwtMiddleware.Process, r.RoleMiddleware.Process(models.RoleAdmin))
	{
		adminCalendar.POST("/add-close-interval", r.AdminCalendarHandler.AddCloseDateInterval)
		adminCalendar.POST("/remove-close-interval", r.AdminCalendarHandler.RemoveCloseDateInterval)
		adminCalendar.GET("/get-visit-stats-list", r.AdminCalendarHandler.GetVisitStatsList)
		adminCalendar.POST("/toggle-visit", r.AdminCalendarHandler.ToggleVisit)
	}

	userCalendar := rg.Group("/user/calendar")
	userCalendar.Use(r.JwtMiddleware.Process)
	{
		userCalendar.GET("/all-user-list", r.RoleMiddleware.Process(models.RoleAdmin), r.AdminCalendarHandler.AllUserList)
		userCalendar.PUT("/update-calendar-entry-status", r.RoleMiddleware.Process(models.RoleAdmin), r.AdminCalendarHandler.ChangeEntryStatus)
	}
}
