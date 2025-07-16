package middlewares

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CurrentUserMiddleware struct {
	UserRepo *repositories.UserRepository
}

func NewCurrentUserMiddleware(db *gorm.DB) *CurrentUserMiddleware {
	userRepo := repositories.NewUserRepository(db)
	return &CurrentUserMiddleware{userRepo}
}

func (middleware *CurrentUserMiddleware) Process(ctx *gin.Context) {
	user, err := middleware.UserRepo.FindByUsername(ctx.GetString("username"))
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
		ctx.Abort()
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}
