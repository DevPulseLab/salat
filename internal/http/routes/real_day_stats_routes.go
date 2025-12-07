package routes

import (
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/handlers"
	"github.com/DevPulseLab/salat/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type RealDayStatsRoutes struct {
	RealDayStatsHandler *handlers.RealDayStatsHandler
	JwtMiddleware       *middlewares.JwtMiddleware
	RoleMiddleware      *middlewares.RoleMiddleware
}

func NewRealDayStatsRoutes(
	realDayStatsHandler *handlers.RealDayStatsHandler,
	jwtMiddleware *middlewares.JwtMiddleware,
	roleMiddleware *middlewares.RoleMiddleware,
) *RealDayStatsRoutes {
	return &RealDayStatsRoutes{
		RealDayStatsHandler: realDayStatsHandler,
		JwtMiddleware:       jwtMiddleware,
		RoleMiddleware:      roleMiddleware,
	}
}

func (r *RealDayStatsRoutes) Setup(rg *gin.RouterGroup) {
	stats := rg.Group("/stats")
	{
		stats.POST("/save-number-of-plates", r.JwtMiddleware.Process, r.RoleMiddleware.Process(models.RoleAdmin), r.RealDayStatsHandler.SaveNumberOfPlatesForDay)
		stats.GET("/get-number-of-plates", r.JwtMiddleware.Process, r.RoleMiddleware.Process(models.RoleAdmin), r.RealDayStatsHandler.GetNumberOfPlatesForDay)
		stats.POST("/increment-number-of-plates", r.RealDayStatsHandler.IncrementNumberOfPlatesForDay)
	}
}
