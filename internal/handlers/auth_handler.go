package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	UserRepo *repositories.UserRepository
	Config   *config.Config
	OAuth2   *oauth2.Config
	Provider *oidc.Provider
}

func NewAuthHandler(db *gorm.DB, config *config.Config) *AuthHandler {
	userRepo := repositories.NewUserRepository(db)

	if !config.Azure.Active {
		return &AuthHandler{userRepo, config, nil, nil}
	}

	provider, err := oidc.NewProvider(context.Background(), "https://login.microsoftonline.com/"+config.Azure.TenantID+"/v2.0")
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.Azure.ClientID,
		ClientSecret: config.Azure.ClientSecret,
		RedirectURL:  config.Azure.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &AuthHandler{userRepo, config, oauth2Config, provider}
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

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered"})
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

	expirationTime := time.Now().Truncate(24 * time.Hour).Add(24*time.Hour - time.Second)
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

func (handler *AuthHandler) SSOLogin(ctx *gin.Context) {
	url := handler.OAuth2.AuthCodeURL("state")
	ctx.Redirect(http.StatusFound, url)
}

func (handler *AuthHandler) SSOCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	token, err := handler.OAuth2.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token in response"})
		return
	}

	expirationTime := time.Now().Truncate(24 * time.Hour).Add(24*time.Hour - time.Second)
	claims := &dto.Claims{
		Username: idToken,
		Role:     models.RoleUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(handler.Config.Jwt.Secret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}
