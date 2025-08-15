package http

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/dbconn"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run(config *config.Config, logger *logrus.Logger) error {
	// Init router
	router := gin.Default()
	InitializeRoutes(router, dbconn.DBSystem, config, logger)

	return router.Run()
}
