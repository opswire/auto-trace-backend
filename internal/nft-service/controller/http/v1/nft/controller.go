package nft

import (
	"car-sell-buy-system/internal/nft-service/domain/nft"
	"car-sell-buy-system/pkg/handler"
	"car-sell-buy-system/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service interface {
	StoreNft(ctx context.Context, dto nft.StoreNftDTO) (nft.Nft, error)
	GetNftByVin(ctx context.Context, vin string) (nft.Nft, error)
	AddServiceRecordByVin(ctx context.Context, vin string, dto nft.AddServiceRecordDTO) (nft.Nft, error)
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
	h := router.Group("/nfts")
	{
		// Protected
		//h.Use(middleware.RequiredAuthMiddleware(ctrl.handler.Logger))
		h.POST("/", ctrl.storeNft)
		h.GET("/:vin", ctrl.getNftByVin)
		h.POST("/:vin/record", ctrl.addServiceRecord)
	}
}

// getNftByVin godoc
//
//	@Summary		Get nft by Vin
//	@Description	Get nft by Car Vin
//	@Tags			Nfts
//	@Accept			json
//	@Produce		json
//	@Param			vin	path		string	true	"Vin"
//	@Success		200	{object}	handler.BasicResponseDTO
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			/api/v1/nfts/{vin} [get]
func (ctrl *Controller) getNftByVin(c *gin.Context) {
	vin, err := ctrl.handler.ParseStringFromPath(c, "vin")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Nft not found. Invalid vin")

		return
	}

	nftToken, err := ctrl.service.GetNftByVin(c.Request.Context(), vin)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Nft not found. Internal error.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newResponse(nftToken),
	})
}

// storeNft godoc
//
//	@Summary		Create new nft
//	@Description	Create new nft
//	@Tags			Nfts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		StoreNftRequest	true	"Nft data"
//	@Success		201		{object}	handler.BasicResponseDTO{data=nft.Response}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Router			/api/v1/nfts [post]
func (ctrl *Controller) storeNft(c *gin.Context) {
	var request StoreNftRequest
	if err := c.ShouldBind(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Nft store error. Invalid request body.")
		return
	}

	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	nftToken, err := ctrl.service.StoreNft(
		ctx,
		request.ToDTO(),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Nft store error. Internal error.")
		return
	}

	ctrl.handler.Logger.Info("Nft with Token ID %d created successfully!", nftToken.TokenId)

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newResponse(nftToken),
	})
}

// addServiceRecord godoc
//
//	@Summary		Add nft service record
//	@Description	Add nft service record
//	@Tags			Nfts
//	@Accept			json
//	@Produce		json
//	@Param			vin	path		string	true	"Vin"
//	@Success		200	{object}	handler.BasicResponseDTO
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Router			/api/v1/nfts/{vin}/record [post]
func (ctrl *Controller) addServiceRecord(c *gin.Context) {
	vin, err := ctrl.handler.ParseStringFromPath(c, "vin")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Nft not found. Invalid vin")
		return
	}

	var request AddServiceRecordRequest
	if err = c.ShouldBind(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Nft store error. Invalid request body.")
		return
	}

	nftToken, err := ctrl.service.AddServiceRecordByVin(
		c.Request.Context(),
		vin,
		request.ToDTO(),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Nft not found. Internal error.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newResponse(nftToken),
	})
}
