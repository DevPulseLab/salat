package task

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/services"
	"github.com/sirupsen/logrus"
	"github.com/uniplaces/carbon"
	"gorm.io/gorm"
)

type CheckReservedRequests struct {
	Config           *config.Config
	MessagingService *services.MessagingService
	CalendarRepo     *repositories.CalendarRepository
	Logger           *logrus.Logger
}

func NewCheckReservedRequests(config *config.Config, db *gorm.DB, logger *logrus.Logger) *CheckReservedRequests {
	ms := services.NewMessagingService(config.Slack.Token, db)
	calendarRepo := repositories.NewCalendarRepository(db)
	return &CheckReservedRequests{config, ms, calendarRepo, logger}
}

func (task *CheckReservedRequests) Execute() {
	now := carbon.Now()
	countReservedRequests, err := task.CalendarRepo.CountReservedByDate(now)
	if err != nil {
		countReservedRequests = 0
		task.Logger.Errorf("Error while count reserved by date: %s", err.Error())
	}

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
