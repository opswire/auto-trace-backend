package v1

import (
	"car-sell-buy-system/config"
	nftController "car-sell-buy-system/internal/nft-service/controller/http/v1/nft"
	"car-sell-buy-system/internal/nft-service/domain/nft"
	"car-sell-buy-system/pkg/logger"
	"github.com/gin-gonic/gin"
)

type V1 struct {
	nftService *nft.Service
	logger     logger.Interface
	config     *config.Config
}

func NewController(
	nftService *nft.Service,
	logger logger.Interface,
	config *config.Config,
) *V1 {
	return &V1{
		nftService: nftService,
		logger:     logger,
		config:     config,
	}
}

func (ctrl *V1) InitAPI(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		nftController.NewController(ctrl.logger, ctrl.nftService).InitAPI(v1)
	}
}
