package task

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SendMessageToChanelTaks struct {
	Config           *config.Config
	MessagingService *service.MessagingService
	log              *logrus.Logger
}

func NewSendMessageToChannelTask(config *config.Config, db *gorm.DB, log *logrus.Logger) *SendMessageToChanelTaks {
	ms := service.NewMessagingService(config.Slack.Token, db)
	return &SendMessageToChanelTaks{config, ms, log}
}

func (task *SendMessageToChanelTaks) Execute() {
	ms := task.MessagingService

	err := ms.PostToChannel(task.Config.Slack.BroadcastChannel, ":green_salad: Bitte denkt an Euren Eintrag in die Satatbar-App (https://salatbar.secova.net/) DANKESCHÖN :green_salad:")
	if err != nil {
		task.log.Errorf("Error while sending message to channel: %s", err.Error())
	}
}
