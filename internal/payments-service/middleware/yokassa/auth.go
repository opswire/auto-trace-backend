package yokassa

import (
	"car-sell-buy-system/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strings"
)

func AuthMiddleware(logger logger.Interface, whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, found := slices.BinarySearch(whitelist, c.ClientIP())

		logger.Info(fmt.Sprintf(
			"Попытка аутентификации вебхука Yokassa IP: %s, Whitelist: %s, Auth: %t",
			c.ClientIP(),
			strings.Join(whitelist, ", "),
			found,
		))

		if !found {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"chat":   "Permission denied!",
			})
			return
		}
	}
}
