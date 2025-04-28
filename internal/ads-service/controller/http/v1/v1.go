package v1

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/ads-service/controller/http/v1/ad"
	"car-sell-buy-system/internal/ads-service/controller/http/v1/chat"
	"car-sell-buy-system/internal/ads-service/controller/http/v1/payment"
	"car-sell-buy-system/pkg/logger"
	"github.com/gin-gonic/gin"
)

type V1 struct {
	adService      ad.Service
	paymentService payment.Service
	chatService    chat.Service
	logger         logger.Interface
	config         *config.Config
}

func NewController(
	adService ad.Service,
	paymentService payment.Service,
	chatService chat.Service,
	logger logger.Interface,
	config *config.Config,
) *V1 {
	return &V1{
		adService:      adService,
		paymentService: paymentService,
		chatService:    chatService,
		logger:         logger,
		config:         config,
	}
}

func (ctrl *V1) InitAPI(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		ad.NewController(ctrl.logger, ctrl.adService).InitAPI(v1)
		payment.NewController(ctrl.logger, ctrl.paymentService, ctrl.config).InitAPI(v1)
		chat.NewController(ctrl.logger, ctrl.chatService).InitAPI(v1)
	}
}
