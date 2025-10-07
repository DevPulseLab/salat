package services

import (
	"fmt"

	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
)

type MessagingService struct {
	client   *slack.Client
	UserRepo *repositories.UserRepository
}

func NewMessagingService(slackToken string, db *gorm.DB) *MessagingService {
	userRepo := repositories.NewUserRepository(db)
	return &MessagingService{
		slack.New(slackToken),
		userRepo,
	}
}

func (ms *MessagingService) SendPrivateMessage(userId uint, content string) error {
	userModel, err := ms.UserRepo.FindById(userId)
	if err != nil {
		return fmt.Errorf("user with id %d not found: %w", userId, err)
	}

	return ms.SendPrivateMessageToEmail(userModel.Username, content)
}

func (ms *MessagingService) SendPrivateMessageToEmail(email string, content string) error {
	user, err := ms.client.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("user with email %s not found: %w", email, err)
	}

	channel, _, _, err := ms.client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{user.ID},
	})
	if err != nil {
		return fmt.Errorf("failed to open DM channel with %s: %w", email, err)
	}

	_, _, err = ms.client.PostMessage(channel.ID, slack.MsgOptionText(content, false))
	if err != nil {
		return fmt.Errorf("failed to send private message to %s: %w", email, err)
	}

	return nil
}

func (ms *MessagingService) PostToChannel(channelName, content string) error {
	channelID, err := ms.findChannelByName(channelName)
	if err != nil {
		return fmt.Errorf("channel #%s not found: %w", channelName, err)
	}

	_, _, err = ms.client.PostMessage(channelID, slack.MsgOptionText(content, false))
	if err != nil {
		return fmt.Errorf("failed to post message to channel #%s: %w", channelName, err)
	}

	return nil
}

func (ms *MessagingService) findChannelByName(channelName string) (string, error) {
	channels, _, err := ms.client.GetConversations(&slack.GetConversationsParameters{
		Types: []string{"public_channel", "private_channel"},
		Limit: 1000,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get channels: %w", err)
	}

	for _, channel := range channels {
		if channel.Name == channelName {
			return channel.ID, nil
		}
	}

	return "", fmt.Errorf("channel not found")
}
