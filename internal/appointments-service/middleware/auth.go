package middleware

import (
	"car-sell-buy-system/pkg/grpc/api/sso_server_v1"
	"car-sell-buy-system/pkg/logger"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

const grpcPort = "50051"

func checkAuth(c *gin.Context, logger logger.Interface) (*sso_server_v1.VerifyTokenResponse, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("sso:%s", grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect SSO: %w", err)
	}
	defer conn.Close() // 1

	client := sso_server_v1.NewSsoV1Client(conn)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
	defer cancel()

	logger.Info(fmt.Sprintf("key::: %v", err))
	response, err := client.VerifyToken(
		ctx,
		&sso_server_v1.VerifyTokenRequest{Token: c.GetHeader("Authorization")},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to request SSO: %w", err)
	}

	return response, err
}

// RequiredAuthMiddleware -.
func RequiredAuthMiddleware(logger logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := checkAuth(c, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("RequiredAuthMiddleware: %v", err))

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Доступ запрещен!"})
			c.Abort()
			return
		}

		if response == nil || !response.Valid {
			logger.Error(fmt.Sprintf("RequiredAuthMiddleware: User is not logined: %v", err))

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Доступ запрещен!"})
			c.Abort()
			return
		}

		logger.Debug(fmt.Sprintf("Response from SSO: %d", response.UserId))

		c.Set("userId", response.UserId)
		c.Next()
	}
}
