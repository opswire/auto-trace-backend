package ad

import "car-sell-buy-system/internal/ads-service/entity"

type StoreRequest struct {
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	Vin           string  `json:"vin"`
	Brand         string  `json:"brand"`
	Model         string  `json:"model"`
	YearOfRelease int64   `json:"year_of_release"`
}

func (r StoreRequest) ToDTO() entity.AdStoreDTO {
	return entity.AdStoreDTO{
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

func (r HandleFavoriteRequest) ToDTO() entity.AdHandleFavoriteDTO {
	return entity.AdHandleFavoriteDTO{
		AdId: r.AdId,
	}
}
