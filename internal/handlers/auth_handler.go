package handlers

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	UserRepo *repositories.UserRepository
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	userRepo := repositories.NewUserRepository(db)
	return &AuthHandler{userRepo}
}

func (handler *AuthHandler) Register(ctx *gin.Context) {
	var form forms.RegisterForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := handler.UserRepo.RegisterUser(form.Username, form.Username); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"message": "User registerd"})
}

func (handler *AuthHandler) Login(ctx *gin.Context) {

}
