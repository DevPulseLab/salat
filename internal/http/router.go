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

	router.Use(middlewares.CORSMiddleware())

	router.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"ping": "pong"})
	})

	router.POST("/api/register", authHandler.Register)
	router.POST("/api/login", authHandler.Login)

	router.GET("/api/users/list", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), userHandler.GetUserList)
	router.GET("/api/user/calendar/all-user-list", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), userCalendarHandler.AllUserList)

	router.POST("/api/user/calendar/add", jwtMiddleware.Process, roleMiddleware.Process(models.RoleUser), userCalendarHandler.Add)
	router.GET("/api/user/calendar/current-user-list", jwtMiddleware.Process, roleMiddleware.Process(models.RoleUser), userCalendarHandler.CurrentUserList)
	router.POST("/api/user/calendar/remove-for-current-user", jwtMiddleware.Process, roleMiddleware.Process(models.RoleUser), userCalendarHandler.RemoveEntryForCurrentUser)
	router.PUT("/api/user/calendar/update-calendar-entry-status", jwtMiddleware.Process, roleMiddleware.Process(models.RoleAdmin), userCalendarHandler.ChangeEntryStatus)
}
