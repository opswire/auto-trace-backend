package v1

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/appointments-service/controller/http/v1/appointment"
	"car-sell-buy-system/pkg/logger"
	"github.com/gin-gonic/gin"
)

type V1 struct {
	appointmentService appointment.Service
	logger             logger.Interface
	config             *config.Config
}

func NewController(
	appointmentService appointment.Service,
	logger logger.Interface,
	config *config.Config,
) *V1 {
	return &V1{
		appointmentService: appointmentService,
		logger:             logger,
		config:             config,
	}
}

func (ctrl *V1) InitAPI(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		appointment.NewController(ctrl.logger, ctrl.appointmentService).InitAPI(v1)
	}
}
