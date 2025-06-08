package appointment

import "time"

type Appointment struct {
	ID          int64
	Start       time.Time
	Duration    int64
	Location    string
	AdId        int64
	SellerId    int64
	BuyerId     int64
	IsConfirmed bool
	IsCanceled  bool
	AdTitle     string
	BuyerName   string
	SellerName  string
	IsBuyer     bool
}
