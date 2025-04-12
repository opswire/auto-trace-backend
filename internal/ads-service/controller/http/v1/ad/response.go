package ad

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/pkg/pagination"
)

type Response struct {
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

func newResponse(ad entity.Ad) Response {
	return Response{
		Id:            ad.Id,
		Title:         ad.Title,
		Description:   ad.Description,
		Price:         ad.Price,
		Vin:           ad.Vin,
		Brand:         ad.Brand,
		Model:         ad.Model,
		YearOfRelease: ad.YearOfRelease,
		IsFavorite:    ad.IsFavorite,
		IsTokenMinted: ad.IsTokenMinted,
	}
}

type ListResponse struct {
	Ads   []Response           `json:"ads"`
	Range pagination.ListRange `json:"range"`
}

func newListResponse(ads []entity.Ad, params pagination.Params, count uint64) ListResponse {
	responses := make([]Response, 0, len(ads))

	for _, ad := range ads {
		responses = append(responses, newResponse(ad))
	}

	return ListResponse{
		Ads:   responses,
		Range: pagination.NewListRange(params, count),
	}
}
