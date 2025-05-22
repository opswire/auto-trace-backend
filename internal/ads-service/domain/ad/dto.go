package ad

import (
	"car-sell-buy-system/pkg/pagination"
	"car-sell-buy-system/pkg/sqlutil"
	"car-sell-buy-system/pkg/storage/local"
)

type StoreDTO struct {
	Title           string
	Description     string
	Price           float64
	Vin             string
	Brand           string
	Model           string
	YearOfRelease   int64
	Image           *local.UploadedFile
	CurrentImageUrl string
	Category        string
	RegNumber       string
	Type            string
	Color           string
	Hp              string
	FullWeight      string
	SoloWeight      string
}

type UpdateDTO struct {
	Title           string
	Description     string
	Price           float64
	Vin             string
	Brand           string
	Model           string
	YearOfRelease   int64
	Image           *local.UploadedFile
	CurrentImageUrl string
	Category        string
	RegNumber       string
	Type            string
	Color           string
	Hp              string
	FullWeight      string
	SoloWeight      string
}

type ListDTO struct {
	Filter     sqlutil.FiltersRequest
	Sort       sqlutil.SortsRequest
	Pagination pagination.Params
}

type HandleFavoriteDTO struct {
	AdId int64
}
