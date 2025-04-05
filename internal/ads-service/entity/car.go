package entity

type Car struct {
	Id            int    `json:"id"`
	Vin           string `json:"vin"`
	IsTokenMinted bool   `json:"is_token_minted"`
	Brand         string `json:"brand"`
	Model         string `json:"model"`
	YearOfRelease int    `json:"year_of_release"`
	ImageUrl      string `json:"image_url"`
}
