package main

import (
	"fmt"
	"os"

	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/cron"
	"github.com/DevPulseLab/salat/internal/db/dbconn"
	"github.com/DevPulseLab/salat/internal/http"
	"github.com/DevPulseLab/salat/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	config := config.New()

	level, err := logrus.ParseLevel(config.ErrorLog.Level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unknown Log-Level in config: %v\n", err)
		os.Exit(1)
	}

	loggerService := service.NewErrorLogger(config.ErrorLog.File, true, level)
	logger := loggerService.GetDefaultLogger()

	dbconn.OpenDB(config.Database.Dsn)
	dbconn.RunMigrate(dbconn.DBSystem)

	cron.Start(config, dbconn.DBSystem, logger)

	if err := http.Run(config, logger); err != nil {
		logger.Fatalf("HTTP-Server-Fehler: %v", err)
	}
}
