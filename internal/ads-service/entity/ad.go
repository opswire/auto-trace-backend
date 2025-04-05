package entity

type Ad struct {
	Id          int     `json:"id"`
	Car         Car     `json:"car"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsFavorite  int8    `json:"is_favorite"`
}
