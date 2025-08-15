package handlers

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/DevPulseLab/salat/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepo         *repositories.UserRepository
	MessagingService *service.MessagingService
	Logger           *logrus.Logger
}

func NewUserHandler(db *gorm.DB, config *config.Config, logger *logrus.Logger) *UserHandler {
	userRepo := repositories.NewUserRepository(db)
	ms := service.NewMessagingService(config.Slack.Token, db)
	return &UserHandler{userRepo, ms, logger}
}

func (handler *UserHandler) GetUserList(ctx *gin.Context) {
	users := handler.UserRepo.GetAllUsers()

	userDtos := []dto.User{}
	for _, user := range users {
		userDto := dto.User{
			Id:          user.ID,
			Username:    user.Username,
			Role:        user.Role,
			PenaltyCard: user.PenaltyCard,
		}
		userDtos = append(userDtos, userDto)
	}

	ctx.JSON(http.StatusOK, gin.H{"users": userDtos})
}

func (handler *UserHandler) SetPenaltyCard(ctx *gin.Context) {
	var form forms.PenaltyCardForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.UserRepo.SetPenaltyCard(form.UserId, form.CardType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not set penalty card"})
		return
	}

	if form.CardType != "" {
		err := handler.MessagingService.SendPrivateMessage(
			form.UserId,
			"Du hast eine Strafkarte bekommen: "+handler.parsePenaltyCard(form.CardType))
		if err != nil {
			handler.Logger.Errorf("Could not send message to user %d: %s", form.UserId, err.Error())
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (handler *UserHandler) GetCurrentUserInfo(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	ctx.JSON(http.StatusOK, gin.H{"penaltyCard": user.PenaltyCard})
}

func (handler *UserHandler) parsePenaltyCard(cardType string) string {
	colorMap := map[string]string{
		"yellow": ":large_yellow_square:",
		"red":    ":large_red_square:",
	}

	parsedType, ok := colorMap[cardType]
	if !ok {
		parsedType = cardType
	}

	return parsedType
}
