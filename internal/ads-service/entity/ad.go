package entity

import (
	"car-sell-buy-system/pkg/pagination"
	"car-sell-buy-system/pkg/sqlutil"
)

type Ad struct {
	Id            int64   `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	Vin           string  `json:"vin"`
	Brand         string  `json:"brand"`
	Model         string  `json:"model"`
	YearOfRelease int64   `json:"year_of_release"`
	IsFavorite    bool    `json:"is_favorite"`
	IsTokenMinted bool    `json:"is_token_minted"`
}

type AdStoreDTO struct {
	Title         string
	Description   string
	Price         float64
	Vin           string
	Brand         string
	Model         string
	YearOfRelease int64
}

type AdListDTO struct {
	Filter     sqlutil.FiltersRequest
	Sort       sqlutil.SortsRequest
	Pagination pagination.Params
}

type AdHandleFavoriteDTO struct {
	AdId int64
}
