package http

import (
	"car-sell-buy-system/config"
	_ "car-sell-buy-system/docs" // Swagger docs.
	v1 "car-sell-buy-system/internal/nft-service/controller/http/v1"
	"car-sell-buy-system/internal/nft-service/domain/nft"
	"car-sell-buy-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter - 1
//
//	@title						Nft Service API
//	@version					1.0
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						http://localhost:8686
//	@BasePath					/api/v1
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func NewRouter(
	handler *gin.Engine,
	logger logger.Interface,
	config *config.Config,
	nftService *nft.Service,
) {
	// Options.
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(CORSMiddleware())

	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_custom_metric",
			Help: "",
		},
		[]string{"total_requests"},
	)
	prometheus.MustRegister(metric)

	h := handler.Group("/api")
	{
		v1.NewController(
			nftService,
			logger,
			config,
		).InitAPI(h)
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
