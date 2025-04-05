package v1

import (
	httpdto "car-sell-buy-system/internal/ads-service/controller/http"
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/internal/ads-service/middleware"
	"car-sell-buy-system/internal/ads-service/usecase"
	"car-sell-buy-system/pkg/logger"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type adRoutes struct {
	uc usecase.Ad
	l  logger.Interface
}

func newAdRoutes(handler *gin.RouterGroup, l logger.Interface, uc usecase.Ad) {
	r := &adRoutes{uc, l}

	h := handler.Group("/ads")
	{
		h.Use(middleware.OptionalAuthMiddleware())
		h.GET("", r.list)
		h.GET("/:adId", r.getById)
		h.POST("", r.store)
		h.GET("/:adId/nftInfo", r.getNftInfo)

		// Protected
		h.Use(middleware.RequiredAuthMiddleware())
		h.POST("/favorite", r.handleFavorite)
	}
}

func (a *adRoutes) getById(c *gin.Context) {
	adId, err := strconv.Atoi(c.Param("adId"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	ad, err := a.uc.GetById(c.Request.Context(), adId)
	if err != nil {
		errorResponse(c, http.StatusNotFound, err.Error())

		return
	}

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   ad,
	})
}

type storeRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Car         car     `json:"car"`
}

type car struct {
	Vin           string `json:"vin"`
	Brand         string `json:"brand"`
	Model         string `json:"model"`
	YearOfRelease int    `json:"year_of_release"`
	ImageUrl      string `json:"image_url"`
}

func (a *adRoutes) store(c *gin.Context) {
	var request storeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	ad, err := a.uc.Store(
		c.Request.Context(),
		entity.Ad{
			Title:       request.Title,
			Description: request.Description,
			Price:       request.Price,
			Car: entity.Car{
				Vin:           request.Car.Vin,
				Brand:         request.Car.Brand,
				Model:         request.Car.Model,
				YearOfRelease: request.Car.YearOfRelease,
				ImageUrl:      request.Car.ImageUrl,
			},
		},
	)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	a.l.Info("Ad with ID %d created successfully!", ad.Id)

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   ad,
	})
}

func (a *adRoutes) list(c *gin.Context) {
	var page, perPage int
	perPage, _ = strconv.Atoi(c.Query("per_page"))
	page, _ = strconv.Atoi(c.Query("page"))

	dt := usecase.BasicListRequestDTO{
		Filter: c.QueryMap("filter"),
		Sort:   c.QueryMap("sort"),
		Pagination: sqlutil.Pagination{
			PerPage: perPage,
			Page:    page,
		},
	}

	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	ads, err := a.uc.List(ctx, dt)
	if err != nil {
		errorResponse(c, http.StatusNotFound, err.Error())

		return
	}

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   ads,
	})
}

type handleRequest struct {
	AdId int `json:"ad_id"`
}

func (a *adRoutes) handleFavorite(c *gin.Context) {
	var request handleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	userId, _ := c.Get("userId")

	err := a.uc.HandleFavorite(
		c.Request.Context(),
		request.AdId,
		int(userId.(int64)),
	)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   gin.H{"message": "success"},
	})
}

func (a *adRoutes) getNftInfo(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("adId"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	ad, err := a.uc.GetTokenInfo(c.Request.Context(), 5552)
	if err != nil {
		errorResponse(c, http.StatusNotFound, err.Error())

		return
	}

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   ad,
	})
}
