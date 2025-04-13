package ad

import (
	"car-sell-buy-system/internal/ads-service/domain/ad"
)

type StoreRequest struct {
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	Vin           string  `json:"vin"`
	Brand         string  `json:"brand"`
	Model         string  `json:"model"`
	YearOfRelease int64   `json:"year_of_release"`
}

func (r StoreRequest) ToDTO() ad.StoreDTO {
	return ad.StoreDTO{
		Title:         r.Title,
		Description:   r.Description,
		Price:         r.Price,
		Vin:           r.Vin,
		Brand:         r.Brand,
		Model:         r.Model,
		YearOfRelease: r.YearOfRelease,
	}
}

type HandleFavoriteRequest struct {
	AdId int64 `json:"ad_id"`
}

func (r HandleFavoriteRequest) ToDTO() ad.HandleFavoriteDTO {
	return ad.HandleFavoriteDTO{
		AdId: r.AdId,
	}
}
