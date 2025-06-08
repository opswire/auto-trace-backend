package middleware

import (
	"car-sell-buy-system/pkg/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequiredAuthMiddleware -.
func RequiredAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Данный способ аутентификации не поддерживается"})
			c.Abort()
			return
		}

		if parts[1] == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Для аутентификации требуется токен"})
			c.Abort()
			return
		}

		claims, err := auth.ParseJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен невалидный"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("userId", claims.ID)
		c.Next()
	}
}
