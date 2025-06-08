package appointment

import (
	"car-sell-buy-system/internal/appointments-service/domain/appointment"
	"car-sell-buy-system/internal/appointments-service/middleware"
	"car-sell-buy-system/pkg/handler"
	"car-sell-buy-system/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Service interface {
	StoreAppointment(ctx context.Context, dto appointment.StoreDTO) (*appointment.Appointment, error)
	CheckTimeConflict(ctx context.Context, dto appointment.CheckTimeConflictDTO) (bool, error)
	GetAllAppointmentsByUserId(ctx context.Context) ([]*appointment.Appointment, error)
	GetAppointmentsByDateRange(ctx context.Context, dto appointment.GetAppointmentsByDateRangeDTO) ([]*appointment.Appointment, error)
	ConfirmAppointment(ctx context.Context, id int64) error
	MarkAppointmentAsCanceled(ctx context.Context, id int64) error
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
	h := router.Group("/appointments")
	{
		// Protected
		h.Use(middleware.RequiredAuthMiddleware(ctrl.handler.Logger))
		h.POST("", ctrl.storeAppointment)
		h.GET("", ctrl.listAppointmentsByUserId)
		h.PATCH("/:appId/confirm", ctrl.confirmAppointment)
		h.PATCH("/:appId/cancel", ctrl.cancelAppointment)
	}
}

// Store Appointment
//
//	@Summary		Create new appointment
//	@Description	Create new appointment
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		StoreAppointmentRequest	true	"Appointment data"
//	@Success		201		{object}	handler.BasicResponseDTO{data=appointment.Response}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Router			/api/v1/appointments [post]
func (ctrl *Controller) storeAppointment(c *gin.Context) {
	var request StoreAppointmentRequest
	if err := c.ShouldBind(&request); err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusBadRequest, err, "Appointment store error: "+err.Error())
		return
	}

	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	app, err := ctrl.service.StoreAppointment(
		ctx,
		request.ToDTO(),
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Appointment store error. "+err.Error())
		return
	}

	ctrl.handler.Logger.Info("Appointment with ID %d created successfully!", app.ID)

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newResponse(app),
	})
}

// List Appointments by User ID
//
//	@Summary		List Appointments by User ID
//	@Description	List Appointments by User ID
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BasicResponseDTO{data=appointment.ListResponse}
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		404	{object}	handler.ErrorResponse
//	@Router			/api/v1/appointments [get]
func (ctrl *Controller) listAppointmentsByUserId(c *gin.Context) {
	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	apps, err := ctrl.service.GetAllAppointmentsByUserId(ctx)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusNotFound, err, "Appointments not found.")

		return
	}

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
		Data:   newListResponse(apps),
	})
}

// Confirm Appointment
//
//	@Summary		Confirm appointment
//	@Description	Confirm appointment
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201		{object}	handler.BasicResponseDTO{data=appointment.Response}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Router			/api/v1/appointments/{appId}/confirm [post]
func (ctrl *Controller) confirmAppointment(c *gin.Context) {
	id, err := ctrl.handler.ParseIDFromPath(c, "appId")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Appointment confirm error. Invalid path id.")
		return
	}

	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	err = ctrl.service.ConfirmAppointment(
		ctx,
		id,
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, err.Error())
		return
	}

	ctrl.handler.Logger.Info("Appointment with ID %d confirmed successfully!", id)

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
	})
}

// Cancel Appointment
//
//	@Summary		Cancel appointment
//	@Description	Cancel appointment
//	@Tags			Appointments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201		{object}	handler.BasicResponseDTO{data=appointment.Response}
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Router			/api/v1/appointments/{appId}/cancel [post]
func (ctrl *Controller) cancelAppointment(c *gin.Context) {
	id, err := ctrl.handler.ParseIDFromPath(c, "appId")
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Appointment cancel error. Invalid path id.")
		return
	}

	userId, _ := c.Get("userId")
	ctx := context.WithValue(c.Request.Context(), "userId", userId)

	err = ctrl.service.MarkAppointmentAsCanceled(
		ctx,
		id,
	)
	if err != nil {
		ctrl.handler.ErrorResponse(c, http.StatusInternalServerError, err, "Appointment cancel error. "+err.Error())
		return
	}

	ctrl.handler.Logger.Info("Appointment with ID %d canceled successfully!", id)

	c.JSON(http.StatusOK, handler.BasicResponseDTO{
		Status: http.StatusOK,
	})
}
