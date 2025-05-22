package payment

import (
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/ads-service/domain/payment"
	"car-sell-buy-system/internal/ads-service/middleware"
	"car-sell-buy-system/pkg/handler"
	"car-sell-buy-system/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service interface {
	CreatePayment(ctx context.Context, dto payment.CreatePaymentDto) (payment.Payment, error)
	ProcessWebhook(ctx context.Context, dto payment.ProcessWebhookPaymentDto) (string, error)
}

type Controller struct {
	handler *handler.BaseHandler
	service Service
	config  *config.Config
}

func NewController(l logger.Interface, service Service, config *config.Config) *Controller {
	return &Controller{
		handler.NewBaseHandler(l),
		service,
		config,
	}
}

func (ctrl *Controller) InitAPI(router *gin.RouterGroup) {
	h := router.Group("/payments")
	{
		//h.Use(yokassa.AuthMiddleware(
		//	ctrl.handler.Logger,
		//	strings.Split(ctrl.config.Yokassa.WebhookAllowedIpAddresses, ", "),
		//))
		h.POST("/webhook", ctrl.processWebhook)

		// Protected
		h.Use(middleware.RequiredAuthMiddleware(ctrl.handler.Logger))
		h.POST("", ctrl.createPayment)
	}
}

// Регистрация платежа
//
//	@Summary		Создание платежа
//	@Description	Регистрация платежа и генерация ссылки
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	handler.BasicResponseDTO{data=payment.Response}
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			/api/v1/payments [post]
func (ctrl *Controller) createPayment(c *gin.Context) {
	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	var request CreatePaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Create payment error. Invalid request body.")
		return
	}

	p, err := ctrl.service.CreatePayment(
		ctx,
		request.toDTO(userId.(int64)),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Create payment error. Internal error.")
		return
	}

	jsonPayment, err := json.Marshal(p)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Create json payment error. Internal error.")
		return
	}

	ctrl.handler.Logger.Info(fmt.Sprintf("Payment created successfully!: %s", jsonPayment))

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newResponse(p),
	})
}

// Вебхук для платежа
//
//	@Summary		Нотификация платежа
//	@Description	Обработка вебхука и смена статуса платежа
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	handler.BasicResponseDTO{data=payment.StatusChangedResponse}
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			/api/v1/payments/webhook [post]
func (ctrl *Controller) processWebhook(c *gin.Context) {
	var request WebhookPaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Process webhook payment error. Invalid request body.")
		return
	}

	ctrl.handler.Logger.Info(fmt.Sprintf("Начата обработка вебхука: %v", request.toDTO()))
	newStatus, err := ctrl.service.ProcessWebhook(
		c.Request.Context(),
		request.toDTO(),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Process webhook payment error. Internal error.")
		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   StatusChangedResponse{Status: newStatus},
	})
}
