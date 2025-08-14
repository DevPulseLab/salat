package cron

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/cron/task"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var c *cron.Cron

func Start(config *config.Config, db *gorm.DB, log *logrus.Logger) {
	c = cron.New()

	c.AddFunc("0 13 * * 4", task.NewSendMessageToChannelTask(config, db, log).Execute)
	c.AddFunc("0 9 * * 5", task.NewSendMessageToChannelTask(config, db, log).Execute)

	c.AddFunc("0 8-16 * * *", task.NewCheckReservedRequests(config, db, log).Execute)

	c.Start()
}
