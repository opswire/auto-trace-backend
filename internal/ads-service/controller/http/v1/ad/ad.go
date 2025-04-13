package ad

import (
	"car-sell-buy-system/internal/ads-service/domain/ad"
	"car-sell-buy-system/internal/ads-service/domain/nft"
	"car-sell-buy-system/internal/ads-service/middleware"
	"car-sell-buy-system/pkg/handler"
	"car-sell-buy-system/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service interface {
	List(ctx context.Context, dto ad.ListDTO) ([]ad.Ad, uint64, error)
	GetById(ctx context.Context, id int64) (ad.Ad, error)
	Store(ctx context.Context, dto ad.StoreDTO) (ad.Ad, error)
	HandleFavorite(ctx context.Context, adId, userId int64) error
	GetTokenInfo(ctx context.Context, tokenId int64) (nft.NFT, error)
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
	h := router.Group("/ads")
	{
		h.Use(middleware.OptionalAuthMiddleware(ctrl.handler.Logger))
		h.GET("", ctrl.list)
		h.GET("/:adId", ctrl.getById)
		h.POST("", ctrl.store)
		h.GET("/:adId/nftInfo", ctrl.getNftInfo)

		// Protected
		h.Use(middleware.RequiredAuthMiddleware(ctrl.handler.Logger))
		h.POST("/favorite", ctrl.handleFavorite)
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
func (ctrl *Controller) getById(c *gin.Context) {
	adId, err := ctrl.handler.ParseIDFromPath(c, "adId")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad not found. Invalid id")

		return
	}

	adv, err := ctrl.service.GetById(c.Request.Context(), adId)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Ad not found. Internal error.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newResponse(adv),
	})
}

// Store Ad
//
//	@Summary		Create new advertisement
//	@Description	Create new car advertisement
//	@Tags			Ads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		StoreRequest	true	"Ad data"
//	@Success		201		{object}	handler.BasicResponseDTO{data=ad.Response}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Router			/api/v1/ads [post]
func (ctrl *Controller) store(c *gin.Context) {
	var request StoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad store error. Invalid request body.")
		return
	}

	adv, err := ctrl.service.Store(
		c.Request.Context(),
		request.ToDTO(),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Ad store error. Internal error.")
		return
	}

	ctrl.handler.Logger.Info("Ad with ID %d created successfully!", adv.Id)

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newResponse(adv),
	})
}

// List Ads
//
//	@Summary		Get ads list
//	@Description	Get paginated and filtered list of ads
//	@Tags			Ads
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int		false	"Page number"		default(1)
//	@Param			limit	query		int		false	"Items per page"	default(10)
//	@Param			filter	query		string	false	"Filter criteria (key=value)"
//	@Param			sort	query		string	false	"Sort field and direction (field=asc|desc)"
//	@Success		200		{object}	handler.BasicResponseDTO{data=ad.ListResponse}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		404		{object}	handler.ErrorResponse
//	@Router			/api/v1/ads [get]
func (ctrl *Controller) list(c *gin.Context) {
	paginationParams, err := ctrl.handler.ParsePaginationParams(c)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Pagination params is not valid.")
	}

	dto := ad.ListDTO{
		Filter:     c.QueryMap("filter"),
		Sort:       c.QueryMap("sort"),
		Pagination: paginationParams,
	}

	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	ads, count, err := ctrl.service.List(ctx, dto)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Ads not found.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newListResponse(ads, paginationParams, count),
	})
}

// Handle Favorite
//
//	@Summary		Add/remove ad to favorites
//	@Description	Toggle ad in user's favorites
//	@Tags			Ads
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		HandleFavoriteRequest	true	"Ad ID"
//	@Success		200		{object}	handler.BasicResponseDTO
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Router			/api/v1/favorites [post]
func (ctrl *Controller) handleFavorite(c *gin.Context) {
	var request HandleFavoriteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad favorite error. Invalid request body.")

		return
	}

	userId, _ := c.Get("userId")

	err := ctrl.service.HandleFavorite(
		c.Request.Context(),
		request.AdId,
		userId.(int64),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Ad favorite error. Internal error.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   gin.H{"message": "success"},
	})
}

// Get NFT Info
//
//	@Summary		Get NFT metadata
//	@Description	Get NFT information for ad
//	@Tags			NFT
//	@Accept			json
//	@Produce		json
//	@Param			adId	path		int	true	"Ad ID"
//	@Success		200		{object}	handler.BasicResponseDTO{data=nft.NFT}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		404		{object}	handler.ErrorResponse
//	@Router			/api/v1/ads/{adId}/nft [get]
func (ctrl *Controller) getNftInfo(c *gin.Context) {
	_, err := ctrl.handler.ParseIDFromPath(c, "adId")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad getNftInfo error. Invalid request body.")

		return
	}

	adv, err := ctrl.service.GetTokenInfo(c.Request.Context(), 5552)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Ad getNftInfo error. Internal error.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   adv,
	})
}
