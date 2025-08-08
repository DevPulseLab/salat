package cron

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/cron/task"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var c *cron.Cron

func Start(config *config.Config, db *gorm.DB) {
	c = cron.New()

	c.AddFunc("0 13 * * 4", task.NewSendMessageToChanelTaks(config, db).Execute)
	c.AddFunc("0 9 * * 5", task.NewSendMessageToChanelTaks(config, db).Execute)

	c.AddFunc("0 8-16 * * *", task.NewCheckReservedRequests(config, db).Execute)

	c.Start()
}
