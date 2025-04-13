package ad

import (
	"car-sell-buy-system/pkg/pagination"
	"car-sell-buy-system/pkg/sqlutil"
)

type StoreDTO struct {
	Title         string
	Description   string
	Price         float64
	Vin           string
	Brand         string
	Model         string
	YearOfRelease int64
}

type ListDTO struct {
	Filter     sqlutil.FiltersRequest
	Sort       sqlutil.SortsRequest
	Pagination pagination.Params
}

type HandleFavoriteDTO struct {
	AdId int64
}
