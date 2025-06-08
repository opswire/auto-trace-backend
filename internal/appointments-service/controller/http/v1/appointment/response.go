package appointment

import (
	"car-sell-buy-system/internal/appointments-service/domain/appointment"
	"time"
)

type Response struct {
	ID          int64     `json:"id"`
	Start       time.Time `json:"start"`
	Duration    int64     `json:"duration"`
	Location    string    `json:"location"`
	AdId        int64     `json:"ad_id"`
	SellerId    int64     `json:"seller_id"`
	BuyerId     int64     `json:"buyer_id"`
	IsConfirmed bool      `json:"is_confirmed"`
	IsCanceled  bool      `json:"is_canceled"`
	AdTitle     string    `json:"ad_title"`
	BuyerName   string    `json:"buyer_name"`
	SellerName  string    `json:"seller_name"`
}

func newResponse(app *appointment.Appointment) Response {
	return Response{
		ID:          app.ID,
		Start:       app.Start,
		Duration:    app.Duration,
		Location:    app.Location,
		AdId:        app.AdId,
		SellerId:    app.SellerId,
		BuyerId:     app.BuyerId,
		IsConfirmed: app.IsConfirmed,
		IsCanceled:  app.IsCanceled,
		AdTitle:     app.AdTitle,
		BuyerName:   app.BuyerName,
		SellerName:  app.SellerName,
	}
}

type ListResponse struct {
	Appointments []Response `json:"appointments"`
}

func newListResponse(apps []*appointment.Appointment) ListResponse {
	responses := make([]Response, 0, len(apps))

	for _, app := range apps {
		responses = append(responses, newResponse(app))
	}

	return ListResponse{
		Appointments: responses,
	}
}
