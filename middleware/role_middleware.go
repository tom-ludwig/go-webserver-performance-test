package middleware

import (
	"go-webserver-performance-test/models"
	"go-webserver-performance-test/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: Check if its possible to inject the decoded token into the gin handler func tree
func RoleMiddleware(requiredRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		claims, err := utils.DecodeToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.Next()
	}
}
