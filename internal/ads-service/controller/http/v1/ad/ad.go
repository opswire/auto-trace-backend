package ad

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/internal/ads-service/middleware"
	"car-sell-buy-system/internal/ads-service/usecase"
	"car-sell-buy-system/pkg/handler"
	"car-sell-buy-system/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	handler *handler.BaseHandler
	uc      usecase.Ad
}

func InitAdRoutes(ginHandler *gin.RouterGroup, l logger.Interface, uc usecase.Ad) {
	r := &Controller{handler.NewBaseHandler(l), uc}

	h := ginHandler.Group("/ads")
	{
		h.Use(middleware.OptionalAuthMiddleware(l))
		h.GET("", r.list)
		h.GET("/:adId", r.getById)
		h.POST("", r.store)
		h.GET("/:adId/nftInfo", r.getNftInfo)

		// Protected
		h.Use(middleware.RequiredAuthMiddleware(l))
		h.POST("/favorite", r.handleFavorite)
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
//	@Router			/ads/{id} [get]
func (ctrl *Controller) getById(c *gin.Context) {
	adId, err := ctrl.handler.ParseIDFromPath(c, "adId")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad not found. Invalid id")

		return
	}

	ad, err := ctrl.uc.GetById(c.Request.Context(), adId)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Ad not found. Internal error.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newResponse(ad),
	})
}

// Store Ad
// @Summary Create new advertisement
// @Description Create new car advertisement
// @Tags Ads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body StoreRequest true "Ad data"
// @Success 201 {object} handler.BasicResponseDTO{data=ad.Response}
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /ads [post]
func (ctrl *Controller) store(c *gin.Context) {
	var request StoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad store error. Invalid request body.")
		return
	}

	ad, err := ctrl.uc.Store(
		c.Request.Context(),
		request.ToDTO(),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Ad store error. Internal error.")
		return
	}

	ctrl.handler.Logger.Info("Ad with ID %d created successfully!", ad.Id)

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newResponse(ad),
	})
}

// List Ads
// @Summary Get ads list
// @Description Get paginated and filtered list of ads
// @Tags Ads
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param filter query string false "Filter criteria (key=value)"
// @Param sort query string false "Sort field and direction (field=asc|desc)"
// @Success 200 {object} handler.BasicResponseDTO{data=ad.ListResponse}
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Router /ads [get]
func (ctrl *Controller) list(c *gin.Context) {
	paginationParams, err := ctrl.handler.ParsePaginationParams(c)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Pagination params is not valid.")
	}

	dto := entity.AdListDTO{
		Filter:     c.QueryMap("filter"),
		Sort:       c.QueryMap("sort"),
		Pagination: paginationParams,
	}

	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	ads, count, err := ctrl.uc.List(ctx, dto)
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
// @Summary Add/remove ad to favorites
// @Description Toggle ad in user's favorites
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body HandleFavoriteRequest true "Ad ID"
// @Success 200 {object} handler.BasicResponseDTO
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /ads/favorites [post]
func (ctrl *Controller) handleFavorite(c *gin.Context) {
	var request HandleFavoriteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad favorite error. Invalid request body.")

		return
	}

	userId, _ := c.Get("userId")

	err := ctrl.uc.HandleFavorite(
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
// @Summary Get NFT metadata
// @Description Get NFT information for ad
// @Tags NFT
// @Accept json
// @Produce json
// @Param adId path int true "Ad ID"
// @Success 200 {object} handler.BasicResponseDTO{data=webapi.NftInfo}
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Router /ads/{adId}/nft [get]
func (ctrl *Controller) getNftInfo(c *gin.Context) {
	_, err := ctrl.handler.ParseIDFromPath(c, "adId")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Ad getNftInfo error. Invalid request body.")

		return
	}

	ad, err := ctrl.uc.GetTokenInfo(c.Request.Context(), 5552)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Ad getNftInfo error. Internal error.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   ad,
	})
}
