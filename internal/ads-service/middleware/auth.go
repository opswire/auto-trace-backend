package middleware

import (
	"car-sell-buy-system/pkg/grpc/api/sso_server_v1"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

const grpcPort = 50051

// OptionalAuthMiddleware -.
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := grpc.NewClient("sso:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}
		defer conn.Close() // 1

		client := sso_server_v1.NewSsoV1Client(conn)

		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
		defer cancel()

		response, err := client.VerifyToken(
			ctx,
			&sso_server_v1.VerifyTokenRequest{Token: c.GetHeader("Authorization")},
		)

		if err != nil {
			log.Printf("failed to request: %v\n", err)
		}

		if response == nil || !response.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Доступ запрещен!"})
			c.Abort()
			return
		}

		c.Set("userId", response.UserId)
	}
}

// RequiredAuthMiddleware -.
func RequiredAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := grpc.NewClient("sso:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}
		defer conn.Close() // 1

		client := sso_server_v1.NewSsoV1Client(conn)

		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
		defer cancel()

		response, err := client.VerifyToken(
			ctx,
			&sso_server_v1.VerifyTokenRequest{Token: c.GetHeader("Authorization")},
		)

		if err != nil {
			log.Printf("failed to request: %v\n", err)
		}

		if response == nil || !response.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Доступ запрещен!"})
			c.Abort()
			return
		}

		c.Set("userId", response.UserId)
	}
}
