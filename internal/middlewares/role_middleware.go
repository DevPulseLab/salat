package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleMiddleware struct {
}

func NewRoleMiddleware() *RoleMiddleware {
	return &RoleMiddleware{}
}

func (middleware *RoleMiddleware) Process(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			ctx.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				ctx.Next()
				return
			}
		}

		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
		ctx.Abort()
	}
}
