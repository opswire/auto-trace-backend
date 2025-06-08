package appointment

import (
	"time"
)

type StoreDTO struct {
	Start       time.Time
	Duration    int64
	Location    string
	AdId        int64
	SellerId    int64
	BuyerId     int64
	IsConfirmed bool
	IsSuccess   bool
}

type CheckTimeConflictDTO struct {
	AdId      int64
	StartTime time.Time
	Duration  int64
}

type GetAppointmentsByDateRangeDTO struct {
	AdId      int64
	StartDate time.Time
	EndDate   time.Time
}

//type ListDTO struct {
//	Filter     sqlutil.FiltersRequest
//	Sort       sqlutil.SortsRequest
//	Pagination pagination.Params
//}
