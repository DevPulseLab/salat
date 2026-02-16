package http

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/http/routes"
	"github.com/DevPulseLab/salat/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(
	router *gin.Engine,
	jwtMiddleware *middlewares.JwtMiddleware,
	authRoutes *routes.AuthRoutes,
	userRoutes *routes.UserRoutes,
	userCalendarRoutes *routes.UserCalendarRoutes,
	adminCalendarRoutes *routes.AdminCalendarRoutes,
	realDayStatsRoutes *routes.RealDayStatsRoutes,
) {
	router.Use(middlewares.CORSMiddleware())

	router.StaticFile("/", "public/index.html")
	router.Static("/public", "public")
	router.Static("/assets", "public/assets")

	api := router.Group("/api")
	{
		api.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, map[string]string{"ping": "pong"})
		})

		authRoutes.Setup(api)

		api.Use(jwtMiddleware.Process)

		userRoutes.Setup(api)
		userCalendarRoutes.Setup(api)
		realDayStatsRoutes.Setup(api)
		adminCalendarRoutes.Setup(api)
	}
}
