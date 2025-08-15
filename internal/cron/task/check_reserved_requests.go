package task

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/helper"
	"github.com/DevPulseLab/salat/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/uniplaces/carbon"
	"gorm.io/gorm"
)

type CheckReservedRequests struct {
	Config           *config.Config
	MessagingService *service.MessagingService
	CalendarRepo     *repositories.CalendarRepository
	Logger           *logrus.Logger
}

func NewCheckReservedRequests(config *config.Config, db *gorm.DB, logger *logrus.Logger) *CheckReservedRequests {
	ms := service.NewMessagingService(config.Slack.Token, db)
	calendarRepo := repositories.NewCalendarRepository(db, helper.NewDateHelper())
	return &CheckReservedRequests{config, ms, calendarRepo, logger}
}

func (task *CheckReservedRequests) Execute() {
	now := carbon.Now()
	countReservedRequests := task.CalendarRepo.CountReservedForDate(now)
	task.Logger.Debugf("Count reserved requests: %d", countReservedRequests)

	if countReservedRequests > 0 {
		err := task.MessagingService.SendPrivateMessageToEmail(
			task.Config.Slack.UserAdminEmail,
			"Es gibt noch ungenehmigte Einträge")

		if err != nil {
			task.Logger.Errorf("Error while sending message: %s", err.Error())
		}
	}
}
