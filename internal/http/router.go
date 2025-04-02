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
	userCalendarHandler := handlers.NewUserCalendarHandler(db)
	adminCalendarHandler := handlers.NewAdminCalendarHandler(db)
	realDayStatsHandler := handlers.NewRealDayStatsHandler(db)

	router.Use(middlewares.CORSMiddleware())

	router.StaticFile("/", "public/index.html")
	router.Static("/public", "public")
	router.Static("/assets", "public/assets")

	router.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"ping": "pong"})
	})

	router.POST("/api/register/cloudflare", authHandler.CloudflareSSO)
	router.POST("/api/register", authHandler.Register)
	router.POST("/api/login", authHandler.Login)

	router.GET("/api/users/list", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), userHandler.GetUserList)
	router.GET("/api/user/calendar/all-user-list", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), adminCalendarHandler.AllUserList)

	router.POST("/api/user/calendar/add", jwtMiddleware.Process, roleMiddleware.Process(models.RoleUser), userCalendarHandler.Add)
	router.GET("/api/user/calendar/current-user-list", jwtMiddleware.Process, roleMiddleware.Process(models.RoleUser), userCalendarHandler.CurrentUserList)
	router.POST("/api/user/calendar/remove-for-current-user", jwtMiddleware.Process, roleMiddleware.Process(models.RoleUser), userCalendarHandler.RemoveEntryForCurrentUser)
	router.PUT("/api/user/calendar/update-calendar-entry-status", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), adminCalendarHandler.ChangeEntryStatus)

	router.POST("/api/stats/save-number-of-plates", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), realDayStatsHandler.SaveNumberOfPlatesForDay)
	router.GET("/api/stats/get-number-of-plates", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), realDayStatsHandler.GetNumberOfPlatesForDay)
	router.POST("/api/stats/increment-number-of-plates", realDayStatsHandler.IncrementNumberOfPlatesForDay)

	router.POST("/api/admin/calendar/add-close-interval", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), adminCalendarHandler.AddCloseDateInterval)
	router.POST("/api/admin/calendar/remove-close-interval", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), adminCalendarHandler.RemoveCloseDateInterval)
	router.GET("/api/admin/calendar/get-visit-stats-list", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), adminCalendarHandler.GetVisitStatsList)
	router.POST("/api/admin/calendar/toggle-visit", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), adminCalendarHandler.ToggleVisit)
	router.GET("/api/user/calendar/get-close-intervals", jwtMiddleware.Process, roleMiddleware.Process(models.RoleUser, models.RoleAdmin), userCalendarHandler.GetCloseDateInterval)
}
