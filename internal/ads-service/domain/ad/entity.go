package ad

import (
	"time"
)

type Ad struct {
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
	ChatExists    bool
	Promotion     Promotion
}

type Promotion struct {
	Status    *string
	ExpiresAt *time.Time
	Enabled   *bool
	TariffId  *int
}
