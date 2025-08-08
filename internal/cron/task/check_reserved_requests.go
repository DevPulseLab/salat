package task

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/helper"
	"github.com/DevPulseLab/salat/internal/service"
	"github.com/uniplaces/carbon"
	"gorm.io/gorm"
)

type CheckReservedRequests struct {
	Config           *config.Config
	MessagingService *service.MessagingService
	CalendarRepo     *repositories.CalendarRepository
}

func NewCheckReservedRequests(config *config.Config, db *gorm.DB) *CheckReservedRequests {
	ms := service.NewMessagingService(config.Slack.Token, db)
	calendarRepo := repositories.NewCalendarRepository(db, helper.NewDateHelper())
	return &CheckReservedRequests{config, ms, calendarRepo}
}

func (task *CheckReservedRequests) Execute() {
	now := carbon.Now()
	countReservedRequests := task.CalendarRepo.CountReservedForDate(now)
	if countReservedRequests > 0 {
		task.MessagingService.SendPrivateMessageToEmail(
			task.Config.Slack.UserAdminEmail,
			"Es gibt noch ungenehmigte Einträge")
	}
}
