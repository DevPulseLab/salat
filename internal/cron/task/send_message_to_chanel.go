package task

import (
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/service"
	"gorm.io/gorm"
)

type SendMessageToChanelTaks struct {
	Config           *config.Config
	MessagingService *service.MessagingService
}

func NewSendMessageToChanelTaks(config *config.Config, db *gorm.DB) *SendMessageToChanelTaks {
	ms := service.NewMessagingService(config.Slack.Token, db)
	return &SendMessageToChanelTaks{config, ms}
}

func (task *SendMessageToChanelTaks) Execute() {
	ms := task.MessagingService

	ms.PostToChannel(task.Config.Slack.BroadcastChannel, ":green_salad: Bitte denkt an Euren Eintrag in die Satatbar-App (https://salatbar.secova.net/) DANKESCHÖN :green_salad:")
}
