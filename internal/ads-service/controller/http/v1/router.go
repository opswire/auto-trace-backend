package v1

import (
	_ "car-sell-buy-system/docs" // Swagger docs.
	"car-sell-buy-system/internal/ads-service/controller/http/v1/ad"
	"car-sell-buy-system/internal/ads-service/usecase"
	"car-sell-buy-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
//
//	@title						Ads Service API
//	@version					1.0
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:8989
//	@BasePath					/api/v1
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func NewRouter(handler *gin.Engine, logger logger.Interface, adUseCase usecase.Ad) {
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
		ad.InitAdRoutes(h, logger, adUseCase)
	}
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods",
			"POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
