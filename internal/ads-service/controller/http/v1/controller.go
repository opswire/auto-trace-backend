package v1

import (
	"car-sell-buy-system/internal/ads-service/controller/http/v1/ad"
	"car-sell-buy-system/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	adService ad.Service
	logger    logger.Interface
}

func NewController(adService ad.Service, logger logger.Interface) *Controller {
	return &Controller{
		adService: adService,
		logger:    logger,
	}
}

func (ctrl *Controller) InitAPI(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		ad.NewController(ctrl.logger, ctrl.adService).InitAPI(v1)
	}
}
