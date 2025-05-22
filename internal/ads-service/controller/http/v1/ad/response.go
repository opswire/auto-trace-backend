package ad

import (
	"car-sell-buy-system/internal/ads-service/domain/ad"
	"car-sell-buy-system/pkg/pagination"
	"time"
)

type Response struct {
	Id            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	Vin           string    `json:"vin"`
	Brand         string    `json:"brand"`
	Model         string    `json:"model"`
	YearOfRelease int64     `json:"year_of_release"`
	IsFavorite    bool      `json:"is_favorite"`
	IsTokenMinted bool      `json:"is_token_minted"`
	ImageUrl      string    `json:"image_url"`
	UserId        int64     `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ChatExists    bool      `json:"chat_exists"`
	Promotion     Promotion `json:"promotion"`
	Category      string    `json:"category"`
	RegNumber     string    `json:"reg_number"`
	Type          string    `json:"type"`
	Color         string    `json:"color"`
	Hp            string    `json:"hp"`
	FullWeight    string    `json:"full_weight"`
	SoloWeight    string    `json:"solo_weight"`
}

type Promotion struct {
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	Enabled   bool      `json:"enabled"`
	TariffId  int       `json:"tariff_id"`
}

func newResponse(ad ad.Ad) Response {
	var promotion Promotion
	if ad.Promotion.ExpiresAt != nil {
		promotion = Promotion{
			Status:    *ad.Promotion.Status,
			ExpiresAt: *ad.Promotion.ExpiresAt,
			Enabled:   ad.Promotion.ExpiresAt.After(time.Now()),
			TariffId:  *ad.Promotion.TariffId,
		}
	}

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
		ImageUrl:      ad.ImageUrl,
		UserId:        ad.UserId,
		CreatedAt:     ad.CreatedAt,
		UpdatedAt:     ad.UpdatedAt,
		ChatExists:    ad.ChatExists,
		Promotion:     promotion,
		Category:      ad.Category,
		RegNumber:     ad.RegNumber,
		Type:          ad.Type,
		Color:         ad.Color,
		Hp:            ad.Hp,
		FullWeight:    ad.FullWeight,
		SoloWeight:    ad.SoloWeight,
	}
}

type ListResponse struct {
	Ads   []Response           `json:"ads"`
	Range pagination.ListRange `json:"range"`
}

func newListResponse(ads []ad.Ad, params pagination.Params, count uint64) ListResponse {
	responses := make([]Response, 0, len(ads))

	for _, adv := range ads {
		responses = append(responses, newResponse(adv))
	}

	return ListResponse{
		Ads:   responses,
		Range: pagination.NewListRange(params, count),
	}
}
