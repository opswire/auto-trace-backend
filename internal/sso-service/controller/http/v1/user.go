package v1

import (
	httpdto "car-sell-buy-system/internal/sso-service/controller/http"
	"car-sell-buy-system/internal/sso-service/entity"
	"car-sell-buy-system/internal/sso-service/middleware"
	"car-sell-buy-system/internal/sso-service/usecase"
	"car-sell-buy-system/pkg/auth"
	"car-sell-buy-system/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type userRoutes struct {
	uc usecase.User
	l  logger.Interface
}

func newAdRoutes(handler *gin.RouterGroup, l logger.Interface, uc usecase.User) {
	r := &userRoutes{uc, l}

	h := handler.Group("/users")
	{
		h.POST("/register", r.register)
		h.POST("/login", r.login)

		// protected
		h.Use(middleware.RequiredAuthMiddleware())
		h.GET("/profile", r.getProfile)
		h.GET("", r.listUsers)
		h.PATCH("/:id/active", r.handleActive)
	}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

func (a *userRoutes) register(c *gin.Context) {
	var request registerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	ad, err := a.uc.Register(
		c.Request.Context(),
		entity.User{
			Email:    request.Email,
			Password: request.Password,
			Name:     request.Name,
		},
	)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	a.l.Info("User with ID %d created successfully!", ad.Id)

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   ad,
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *userRoutes) login(c *gin.Context) {
	var login loginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	user, err := a.uc.GetByEmail(c.Request.Context(), login.Email)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, httpdto.BasicResponseDTO{
			Status: http.StatusNotFound,
			Data:   "пользователь с таким email не найден",
		})
		return
	}

	if err = user.ComparePasswords(login.Password); err != nil {
		c.JSON(http.StatusBadRequest, httpdto.BasicResponseDTO{
			Status: http.StatusBadRequest,
			Data:   "Пароль не совпадает",
		})
		return
	}

	token, err := auth.GenerateJWT(strconv.Itoa(user.Id), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, httpdto.BasicResponseDTO{
			Status: http.StatusBadRequest,
			Data:   "Пользователь не активен",
		})
		return
	}

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data: gin.H{
			"access_token": "Bearer " + token,
		},
	})
}

func (a *userRoutes) getProfile(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email not found"})
		return
	}

	user, err := a.uc.GetByEmail(c.Request.Context(), email.(string))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (a *userRoutes) listUsers(c *gin.Context) {
	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	users, err := a.uc.List(ctx)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   users,
	})
}

func (a *userRoutes) handleActive(c *gin.Context) {
	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	idParam := c.Param("id")
	if idParam == "" {
		errorResponse(c, http.StatusBadRequest, "ID must be integer")

		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if id <= 0 || err != nil {
		errorResponse(c, http.StatusBadRequest, "ID must be integer")

		return
	}

	err = a.uc.HandleActive(ctx, id)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, httpdto.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   "OK",
	})
}
