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
	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepo         *repositories.UserRepository
	MessagingService *service.MessagingService
}

func NewUserHandler(db *gorm.DB, config *config.Config) *UserHandler {
	userRepo := repositories.NewUserRepository(db)
	ms := service.NewMessagingService(config.Slack.Token, db)
	return &UserHandler{userRepo, ms}
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
		handler.MessagingService.SendPrivateMessage(form.UserId, "Du hast eine Strafkarte bekommen: "+form.CardType)
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (handler *UserHandler) GetCurrentUserInfo(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	ctx.JSON(http.StatusOK, gin.H{"penaltyCard": user.PenaltyCard})
}
