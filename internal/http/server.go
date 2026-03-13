package http

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/dbconn"
	"github.com/DevPulseLab/salat/internal/handlers"
	"github.com/DevPulseLab/salat/internal/http/routes"
	"github.com/DevPulseLab/salat/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run(config *config.Config, logger *logrus.Logger) error {
	jwtMiddleware := middlewares.NewJwtMiddleware(config)
	roleMiddleware := middlewares.NewRoleMiddleware()
	currentUserMiddleware := middlewares.NewCurrentUserMiddleware(dbconn.DBSystem)

	authHandler := handlers.NewAuthHandler(dbconn.DBSystem, config)
	userHandler := handlers.NewUserHandler(dbconn.DBSystem, config, logger)
	userCalendarHandler := handlers.NewUserCalendarHandler(dbconn.DBSystem)
	adminCalendarHandler := handlers.NewAdminCalendarHandler(dbconn.DBSystem, config, logger)
	realDayStatsHandler := handlers.NewRealDayStatsHandler(dbconn.DBSystem)

	authRoutes := routes.NewAuthRoutes(authHandler)
	userRoutes := routes.NewUserRoutes(userHandler, roleMiddleware, currentUserMiddleware)
	userCalendarRoutes := routes.NewUserCalendarRoutes(userCalendarHandler, jwtMiddleware, roleMiddleware, currentUserMiddleware)
	adminCalendarRoutes := routes.NewAdminCalendarRoutes(adminCalendarHandler, jwtMiddleware, roleMiddleware)
	realDayStatsRoutes := routes.NewRealDayStatsRoutes(realDayStatsHandler, jwtMiddleware, roleMiddleware)

	router := gin.Default()
	InitializeRoutes(
		router,
		jwtMiddleware,
		authRoutes,
		userRoutes,
		userCalendarRoutes,
		adminCalendarRoutes,
		realDayStatsRoutes,
	)

	return router.Run()
}
