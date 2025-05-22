package chat

import (
	"car-sell-buy-system/internal/ads-service/domain/chat"
	"car-sell-buy-system/internal/ads-service/middleware"
	"car-sell-buy-system/pkg/handler"
	"car-sell-buy-system/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service interface {
	StoreMessage(ctx context.Context, chatId int64, dto chat.StoreMessageDTO) (chat.Message, error)
	StoreChat(ctx context.Context, dto chat.StoreChatDTO) (chat.Chat, error)
	ListChats(ctx context.Context) ([]chat.Chat, int64, error)
	ListMessagesByChatId(ctx context.Context, chatId int64) ([]chat.Message, int64, error)
}

type Controller struct {
	handler *handler.BaseHandler
	service Service
}

func NewController(l logger.Interface, service Service) *Controller {
	return &Controller{
		handler.NewBaseHandler(l),
		service,
	}
}

func (ctrl *Controller) InitAPI(router *gin.RouterGroup) {
	h := router.Group("/chats")
	{
		// Protected
		h.Use(middleware.RequiredAuthMiddleware(ctrl.handler.Logger))
		h.GET("", ctrl.listChats)
		h.POST("", ctrl.storeChat)
		h.GET("/:chatId/messages", ctrl.listMessagesByChatId)
		h.POST("/:chatId/messages", ctrl.storeMessage)
	}
}

// getById godoc
//
//	@Summary		Get advertisement by ID
//	@Description	Get car advertisement details
//	@Tags			Ads
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Ad ID"
//	@Success		200	{object}	handler.BasicResponseDTO
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			/api/v1/ads/{id} [get]
//func (ctrl *Controller) getById(c *gin.Context) {
//	adId, err := ctrl.handler.ParseIDFromPath(c, "adId")
//	if err != nil {
//		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad not found. Invalid id")
//
//		return
//	}
//
//	adv, err := ctrl.service.GetById(c.Request.Context(), adId)
//	if err != nil {
//		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Ad not found. Internal error.")
//
//		return
//	}
//
//	c.JSON(http.StatusOK, handler.BasicResponseDTO{
//		Status: http.StatusOK,
//		Data:   newResponse(adv),
//	})
//}

// Store Chat
//
//	@Summary		Create new chat
//	@Description	Create new chat
//	@Tags			Chats
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		StoreChatRequest	true	"Chat data"
//	@Success		201		{object}	handler.BasicResponseDTO{data=chat.ChatResponse}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Router			/api/v1/chats [post]
func (ctrl *Controller) storeChat(c *gin.Context) {
	var request StoreChatRequest
	if err := c.ShouldBind(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Chat store error. Invalid request body.")
		return
	}

	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	cht, err := ctrl.service.StoreChat(
		ctx,
		request.ToDTO(),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Chat store error. Internal error.")
		return
	}

	ctrl.handler.Logger.Info("Chat with ID %d created successfully!", cht.Id)

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newChatResponse(cht),
	})
}

// Store Message
//
//	@Summary		Create new message
//	@Description	Create new message
//	@Tags			Chats
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		StoreMessageRequest	true	"Message data"
//	@Success		201		{object}	handler.BasicResponseDTO{data=chat.MessageResponse}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Router			/api/v1/chats/{id}/message [post]
func (ctrl *Controller) storeMessage(c *gin.Context) {
	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	var request StoreMessageRequest
	if err := c.ShouldBind(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Message store error. Invalid request body.")
		return
	}

	chatId, err := ctrl.handler.ParseIDFromPath(c, "chatId")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Message store error. Invalid parse chat id.")
		return
	}

	msg, err := ctrl.service.StoreMessage(
		ctx,
		chatId,
		request.ToDTO(),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Message store error. Internal error.")
		return
	}

	ctrl.handler.Logger.Info("Message with ID %d created successfully!", msg.Id)

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newMessageResponse(msg),
	})
}

// List Messages
//
//	@Summary		Get messages list
//	@Description	Get messages list
//	@Tags			Chats
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BasicResponseDTO{data=chat.ListMessageResponse}
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		404	{object}	handler.ErrorResponse
//	@Router			/api/v1/chats/{id}/messages [get]
func (ctrl *Controller) listMessagesByChatId(c *gin.Context) {
	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	chatId, err := ctrl.handler.ParseIDFromPath(c, "chatId")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad not found. Invalid id")

		return
	}

	messages, _, err := ctrl.service.ListMessagesByChatId(ctx, chatId)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Chats not found.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newListMessageResponse(messages),
	})
}

// List Chats
//
//	@Summary		Get chats list
//	@Description	Get chats list
//	@Tags			Chats
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BasicResponseDTO{data=chat.ListChatResponse}
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		404	{object}	handler.ErrorResponse
//	@Router			/api/v1/chats [get]
func (ctrl *Controller) listChats(c *gin.Context) {
	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	chats, _, err := ctrl.service.ListChats(ctx)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Chats not found.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newListChatsResponse(chats),
	})
}

// Delete Ad
//
//	@Summary		Delete advertisement
//	@Description	Delete car advertisement
//	@Tags			Ads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		StoreRequest	true	"Ad data"
//	@Success		201		{object}	handler.BasicResponseDTO{data=string}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Router			/api/v1/ads/{id} [delete]
//func (ctrl *Controller) delete(c *gin.Context) {
//	id, err := ctrl.handler.ParseIDFromPath(c, "adId")
//	if err != nil {
//		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Ad update error. Invalid path id.")
//		return
//	}
//
//	err = ctrl.service.Delete(
//		c.Request.Context(),
//		id,
//	)
//	if err != nil {
//		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Ad update error. Internal error.")
//		return
//	}
//
//	ctrl.handler.Logger.Info("Ad with ID %d updated successfully!", id)
//
//	c.JSON(http.StatusOK, handler.BasicResponseDTO{
//		Status: http.StatusOK,
//		Data:   "success",
//	})
//}
