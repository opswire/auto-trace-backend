package v1

import (
	"car-sell-buy-system/internal/ads-service/usecase"
	"car-sell-buy-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *gin.Engine, l logger.Interface, a usecase.Ad) {
	// Options.
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(CORSMiddleware())

	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_custom_metric",
			Help: ""},
		[]string{"total_requests"},
	)
	prometheus.MustRegister(metric)

	h := handler.Group("/v1")
	{
		newAdRoutes(h, l, a)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
