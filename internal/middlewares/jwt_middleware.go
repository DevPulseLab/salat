package middlewares

import (
	"net/http"
	"strings"

	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	RoleGuest = "guest"
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type JwtMiddleware struct {
	Config *config.Config
}

func NewJwtMiddleware(config *config.Config) *JwtMiddleware {
	return &JwtMiddleware{config}
}

func (middleware *JwtMiddleware) Process(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		ctx.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &dto.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(middleware.Config.Jwt.Secret), nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	// Pass the user information to the next handler
	ctx.Set("username", claims.Username)
	ctx.Set("role", claims.Role)
	ctx.Next()
}
