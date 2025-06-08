package appointment

import (
	"car-sell-buy-system/internal/appointments-service/domain/appointment"
	"time"
)

type StoreAppointmentRequest struct {
	Start    time.Time `json:"start" binding:"required"`
	Duration int64     `json:"duration" binding:"required"`
	Location string    `json:"location" binding:"required"`
	AdId     int64     `json:"ad_id" binding:"required"`
	BuyerId  int64     `json:"buyer_id" binding:"required"`
}

func (r StoreAppointmentRequest) ToDTO() appointment.StoreDTO {
	return appointment.StoreDTO{
		Start:    r.Start,
		Duration: r.Duration,
		Location: r.Location,
		AdId:     r.AdId,
		BuyerId:  r.BuyerId,
	}
}
