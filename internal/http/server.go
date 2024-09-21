package http

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/dbconn"
	"github.com/gin-gonic/gin"
)

func Run() error {
	// Load configuration
	configuration := config.New()

	// Init database
	dbconn.OpenDB(configuration.Database.Dsn)
	dbconn.RunMigrate(dbconn.DBSystem)

	// Init router
	router := gin.Default()
	InitializeRoutes(router, dbconn.DBSystem, configuration)

	return router.Run()
}
