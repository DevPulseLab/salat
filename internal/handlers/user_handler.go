package handlers

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepo *repositories.UserRepository
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	userRepo := repositories.NewUserRepository(db)
	return &UserHandler{userRepo}
}

func (handler *UserHandler) GetUserList(ctx *gin.Context) {
	users := handler.UserRepo.GetAllUsers()

	userDtos := []dto.User{}
	for _, user := range users {
		userDto := dto.User{
			Id:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		}
		userDtos = append(userDtos, userDto)
	}

	ctx.JSON(http.StatusOK, gin.H{"users": userDtos})
}
