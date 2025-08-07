package main

import (
	"log"

	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/cron"
	"github.com/DevPulseLab/salat/internal/db/dbconn"
	"github.com/DevPulseLab/salat/internal/http"
)

func main() {
	// Load configuration
	config := config.New()

	// Init database
	dbconn.OpenDB(config.Database.Dsn)
	dbconn.RunMigrate(dbconn.DBSystem)

	cron.Start(config, dbconn.DBSystem)
	if err := http.Run(config); err != nil {
		log.Fatalln(err)
	}
}
