package ad

import (
	"time"
)

type Ad struct {
	Id            int64
	Title         string
	Description   string
	Price         float64
	Vin           string
	Brand         string
	Model         string
	YearOfRelease int64
	IsFavorite    bool
	IsTokenMinted bool
	ImageUrl      string
	UserId        int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ChatExists    bool
	Promotion     Promotion
	Category      string
	RegNumber     string
	Type          string
	Color         string
	Hp            string
	FullWeight    string
	SoloWeight    string
}

type Promotion struct {
	Status    *string
	ExpiresAt *time.Time
	Enabled   *bool
	TariffId  *int
}
