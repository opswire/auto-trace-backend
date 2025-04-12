package handler

import (
	"car-sell-buy-system/pkg/logger"
	"car-sell-buy-system/pkg/pagination"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	ErrEmptyIDParam        = errors.New("empty id param")
	ErrInvalidID           = errors.New("invalid id param")
	ErrInvalidPerPageParam = errors.New("invalid per_page param")
	ErrInvalidPageParam    = errors.New("invalid page param")
)

type BaseHandler struct {
	Logger logger.Interface
}

func NewBaseHandler(logger logger.Interface) *BaseHandler {
	return &BaseHandler{Logger: logger}
}

type ErrorResponse struct {
	Error string `json:"error" example:"message"`
}

func (h *BaseHandler) ErrorResponse(c *gin.Context, status int, err error, message string) {
	h.Logger.Error("Critical error: %s", err.Error())

	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: message,
	})
}

func (h *BaseHandler) ParseIDFromPath(c *gin.Context, param string) (int64, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return 0, ErrEmptyIDParam
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if id <= 0 || err != nil {
		return 0, ErrInvalidID
	}

	return id, nil
}
func (h *BaseHandler) ParsePaginationParams(c *gin.Context) (pagination.Params, error) {
	var (
		params pagination.Params
		err    error
	)

	perPageParam := c.Query("per_page")
	params.PerPage, err = strconv.ParseUint(perPageParam, 10, 64)
	if err != nil {
		if perPageParam != "" {
			return pagination.Params{}, ErrInvalidPerPageParam
		}
		params.PerPage = pagination.DefaultPerPage
	}

	if params.PerPage < pagination.MinPerPage || params.PerPage > pagination.MaxPerPage {
		return pagination.Params{}, ErrInvalidPerPageParam
	}

	pageParam := c.Query("page")
	params.Page, err = strconv.ParseUint(pageParam, 10, 64)
	if err != nil {
		if pageParam != "" {
			return pagination.Params{}, ErrInvalidPageParam
		}
		params.Page = pagination.DefaultPage
	}

	return params, nil
}
