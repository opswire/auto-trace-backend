package v1

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/payments-service/controller/http/v1/payment"
	"car-sell-buy-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type V1 struct {
	paymentService payment.Service
	logger         logger.Interface
	config         *config.Config
	publisher      *kafka.Writer
}

func NewController(
	paymentService payment.Service,
	logger logger.Interface,
	config *config.Config,
	publisher *kafka.Writer,
) *V1 {
	return &V1{
		paymentService: paymentService,
		logger:         logger,
		config:         config,
		publisher:      publisher,
	}
}

func (ctrl *V1) InitAPI(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		payment.NewController(ctrl.logger, ctrl.paymentService, ctrl.config, ctrl.publisher).InitAPI(v1)
	}
}
