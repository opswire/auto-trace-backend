package v1

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/chats-service/controller/http/v1/chat"
	"car-sell-buy-system/pkg/logger"
	"github.com/gin-gonic/gin"
)

type V1 struct {
	chatService chat.Service
	logger      logger.Interface
	config      *config.Config
}

func NewController(
	chatService chat.Service,
	logger logger.Interface,
	config *config.Config,
) *V1 {
	return &V1{
		chatService: chatService,
		logger:      logger,
		config:      config,
	}
}

func (ctrl *V1) InitAPI(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		chat.NewController(ctrl.logger, ctrl.chatService).InitAPI(v1)
	}
}
