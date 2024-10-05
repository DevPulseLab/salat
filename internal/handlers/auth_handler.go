package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	UserRepo *repositories.UserRepository
	Config   *config.Config
}

func NewAuthHandler(db *gorm.DB, config *config.Config) *AuthHandler {
	userRepo := repositories.NewUserRepository(db)
	return &AuthHandler{userRepo, config}
}

func (handler *AuthHandler) Register(ctx *gin.Context) {
	var form forms.RegisterForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := handler.UserRepo.RegisterUser(form.Username, form.Password, models.RoleUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"message": "User registerd"})
}

func (handler *AuthHandler) Login(ctx *gin.Context) {
	var form forms.LoginForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userRole, err := handler.UserRepo.AuthenticateUser(form.Username, form.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &dto.Claims{
		Username: form.Username,
		Role:     userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(handler.Config.Jwt.Secret)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}
