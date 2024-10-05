package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (handler *UserHandler) GetUserList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Ok", "Role": ctx.GetString("role"), "Username": ctx.GetString("username")})
}
